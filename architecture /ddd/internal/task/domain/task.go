package domain

import (
	"errors"
	"time"
)

// TaskID 用单独类型包装，方便以后替换为 UUID 等
type TaskID string

// TaskStatus 任务状态（DDD 里的值对象）
type TaskStatus string

// 执行状态
const (
	StatusPending   TaskStatus = "PENDING"   // 已创建，等待调度
	StatusRunning   TaskStatus = "RUNNING"   // 执行中
	StatusCompleted TaskStatus = "COMPLETED" // 已完成
	StatusFailed    TaskStatus = "FAILED"    // 失败
)

// Task 聚合根：一个“任务”的领域对象
type Task struct {
	id          TaskID
	name        string
	description string
	status      TaskStatus
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTask 工厂函数：创建任务（领域规则在这里校验）
func NewTask(id TaskID, name, desc string, now time.Time) (*Task, error) {
	if name == "" {
		return nil, errors.New("task name is required")
	}

	// 初始化任务参数
	return &Task{
		id:          id,
		name:        name,
		description: desc,
		status:      StatusPending,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Start 领域行为：开始执行任务
func (t *Task) Start(now time.Time) error {
	if t.status != StatusPending {
		return errors.New("only pending task can be started")
	}
	t.status = StatusRunning
	t.updatedAt = now
	return nil
}

// Complete 领域行为：标记完成
func (t *Task) Complete(now time.Time) error {
	if t.status != StatusRunning {
		return errors.New("only running task can be completed")
	}
	t.status = StatusCompleted
	t.updatedAt = now
	return nil
}

// Fail 领域行为：标记失败
func (t *Task) Fail(now time.Time) error {
	if t.status != StatusRunning && t.status != StatusPending {
		return errors.New("only pending or running task can be failed")
	}
	t.status = StatusFailed
	t.updatedAt = now
	return nil
}

// ID 一些 getter，防止外部随便改内部字段
func (t *Task) ID() TaskID           { return t.id }
func (t *Task) Name() string         { return t.name }
func (t *Task) Description() string  { return t.description }
func (t *Task) Status() TaskStatus   { return t.status }
func (t *Task) CreatedAt() time.Time { return t.createdAt }
func (t *Task) UpdatedAt() time.Time { return t.updatedAt }
