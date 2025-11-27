package dlock

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// init 提供的种子值将默认的 Source 初始化为确定性状态
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 对外可见的一些错误
var (
	ErrAcquireTimeout = errors.New("dlock: acquire lock timeout")
	ErrNotOwner       = errors.New("dlock: lock not owned by this client")
)

// Client 封装 Redis 客户端以及本地说状态
type Client struct {
	rdb *redis.Client

	mu     sync.Mutex
	states map[string]*lockState
}

// 每个 key 对应的本地锁状态实现（实现可重入 + 续约协程管理）
type lockState struct {
	token  string
	count  int
	cancel context.CancelFunc // 用于停止 watchDog
}

// Lock 是用户拿到的锁句柄
type Lock struct {
	key    string
	client *Client
}

// LockOptions 控制加锁行为
type LockOptions struct {
	// Duration 两个瞬间之间的经过时间
	TTL           time.Duration // 锁过期时间（服务端TTL）
	TryTimeout    time.Duration // 总共等待多长时间去获取锁（0 表示一直等到 ctx 取消）
	RetryInterval time.Duration // 每次重试基础间隔
	AutoRenew     bool          // 是否自动续约
	RenewInterval time.Duration // 续约间隔（<= TTL, 默认 TTL/3）
}

// 一些默认值
const (
	defaultTTL           = 10 * time.Second
	defaultRetryInterval = 50 * time.Millisecond
	jitterFactor         = 0.5 // 重试间隔抖动比例
)

// DefaultLockOptions 默认配置
func DefaultLockOptions() LockOptions {
	return LockOptions{
		TTL:           defaultTTL,
		TryTimeout:    0,
		RetryInterval: defaultRetryInterval,
		AutoRenew:     false,
		RenewInterval: 0,
	}
}

// Option 函数式编程
type Option func(*LockOptions)

// WithTTL 初始化 TTL
func WithTTL(ttl time.Duration) Option {
	return func(o *LockOptions) {
		o.TTL = ttl
	}
}

// WithTryTimeout 初始化 TryTimeout
func WithTryTimeout(d time.Duration) Option {
	return func(o *LockOptions) {
		o.TryTimeout = d
	}
}

// WithRetryInterval 初始化 RetryInterval
func WithRetryInterval(d time.Duration) Option {
	return func(o *LockOptions) {
		o.RetryInterval = d
	}
}

// WithAutoRenew 初始化 AutoRenew
func WithAutoRenew(enable bool) Option {
	return func(o *LockOptions) {
		o.AutoRenew = enable
	}
}

// WithRenewInterval 初始化 RenewInterval
func WithRenewInterval(d time.Duration) Option {
	return func(o *LockOptions) {
		o.RenewInterval = d
	}
}

// NewClient 用外部创建好的 go-redis Client 初始化
func NewClient(rdb *redis.Client) *Client {
	return &Client{
		rdb:    rdb,
		states: make(map[string]*lockState),
	}
}

// lua 脚本：解锁（只有 value 匹配时才删除）
var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end
`)

// lua 脚本：续约（只有 value 匹配时才 PEXPIRE）
var renewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
  return 0
end
`)

// TryLock 尝试一次获取锁（非阻塞），成功返回 (*Lock, nil)，失败但没有错误返回 (nil, nil)
func (c *Client) TryLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	// base case
	if ttl < 0 {
		ttl = defaultTTL
	}

	// 先检测本地是否已经持有：实现“可重入锁”
	c.mu.Lock()
	if st, ok := c.states[key]; ok {
		st.count++
		c.mu.Unlock()
		return &Lock{key: key, client: c}, nil
	}
	c.mu.Unlock()

	token := uuid.NewString()
	// setNx 操作执行分布式锁
	ok, err := c.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, err
	}

	// 说明有人持有锁
	if !ok {
		return nil, nil
	}

	// 写入本地状态
	c.mu.Lock()
	c.states[key] = &lockState{
		token: token,
		count: 1,
	}
	c.mu.Unlock()

	return &Lock{key: key, client: c}, nil
}

// Lock 带重试 & 超时获取锁
//
// 语义：
//   - 如果当前进程已经持有该 key 的锁，则只是 count++（可重入），立即返回
//   - 否则循环尝试 SET NX，直到：
//   - 成功拿到锁；或者
//   - TryTimeout 到达；或者
//   - ctx 被取消
func (c *Client) Lock(ctx context.Context, key string, opts ...Option) (*Lock, error) {
	// 默认配置
	cfg := DefaultLockOptions()
	for _, fn := range opts {
		fn(&cfg)
	}
	// base case
	if cfg.TTL <= 0 {
		cfg.TTL = defaultTTL
	}
	if cfg.RetryInterval <= 0 {
		cfg.RetryInterval = defaultRetryInterval
	}

	// 先看看本地是否已经持有（可重入）
	c.mu.Lock()
	if st, ok := c.states[key]; ok {
		st.count++
		c.mu.Unlock()
		return &Lock{key: key, client: c}, nil
	}
	c.mu.Unlock()

	// 计算整体超时 deadline（如果配置了 TryTimeout）
	var deadline time.Time
	if cfg.TryTimeout > 0 {
		// 添加等待时间
		deadline = time.Now().Add(cfg.TryTimeout)
	}

	// 使用 uuid 的方式，创建 token
	token := uuid.NewString()

	for {
		// 先检查 ctx
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 如果有 TryTimeout 限制，也要检查
		if !deadline.IsZero() && time.Now().After(deadline) {
			return nil, ErrAcquireTimeout
		}

		ok, err := c.rdb.SetNX(ctx, key, token, cfg.TTL).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			// 成功拿到锁，记录本地状态 + 启动 watchdog（如果需要）
			var cancel context.CancelFunc
			if cfg.AutoRenew {
				// 启动续约协程
				renewInterval := cfg.RenewInterval
				if renewInterval <= 0 || renewInterval >= cfg.TTL {
					renewInterval = cfg.TTL / 3
					if renewInterval <= 0 {
						renewInterval = time.Millisecond * 100
					}
				}
				var wctx context.Context
				wctx, cancel = context.WithCancel(context.Background())
				go c.watchdog(wctx, key, token, cfg.TTL, renewInterval)
			}

			c.mu.Lock()
			c.states[key] = &lockState{
				token:  token,
				count:  1,
				cancel: cancel,
			}
			c.mu.Unlock()

			return &Lock{key: key, client: c}, nil
		}

		// 没拿到锁：睡一会儿再重试（带一点随机抖动）
		sleep := cfg.RetryInterval
		jitter := time.Duration(rand.Float64() * jitterFactor * float64(sleep))
		sleep = sleep + jitter

		// 防止因 ctx 取消而多睡
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleep):
		}
	}
}

// watchdog 定期续约锁，直到 ctx 被取消或续约失败
func (c *Client) watchdog(ctx context.Context, key, token string, ttl, interval time.Duration) {
	// 为了安全一点，我们续约的 TTL 仍然用原始 ttl
	ttlMs := ttl.Milliseconds()
	if ttlMs <= 0 {
		ttlMs = defaultTTL.Milliseconds()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 调用 Lua 续约脚本；只有 token 匹配时才会 PEXPIRE
			_, err := renewScript.Run(ctx, c.rdb, []string{key}, token, ttlMs).Int()
			if err != nil {
				// 一般可以记录日志，这里我们只选择“停止 watchdog”
				return
			}
		}
	}
}

// Unlock 释放锁。
// 可重入场景下，需要调用 Unlock 与 Lock 调用次数匹配；
// 只有最后一次 Unlock 才会真正删除 Redis 里的 key。
func (l *Lock) Unlock(ctx context.Context) error {
	return l.client.unlock(ctx, l.key)
}

// 内部解锁逻辑
func (c *Client) unlock(ctx context.Context, key string) error {
	c.mu.Lock()
	st, ok := c.states[key]
	if !ok {
		c.mu.Unlock()
		return ErrNotOwner
	}

	st.count--
	if st.count > 0 {
		// 还有重入层数，不真正释放
		c.mu.Unlock()
		return nil
	}

	// 最后一次解锁：停止 watchdog + 删除本地状态
	if st.cancel != nil {
		st.cancel()
	}
	token := st.token
	delete(c.states, key)
	c.mu.Unlock()

	// 用 Lua 脚本安全删除 Redis 锁（防止误删他人锁）
	res, err := unlockScript.Run(ctx, c.rdb, []string{key}, token).Int()
	if err != nil {
		return err
	}
	if res == 0 {
		// 没删掉：要么锁已过期，要么 token 不匹配
		return ErrNotOwner
	}
	return nil
}

// Key 一些辅助方法，方便调试
func (l *Lock) Key() string { return l.key }
