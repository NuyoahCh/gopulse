package main

import (
	"context"
	"fmt"
	dlock "github.com/Nuyoahch/gopulse/lock/distlock"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	client := dlock.NewClient(rdb)

	ctx := context.Background()

	// 获取一把锁，TTL=5s，获取最多等 3s，自动续约
	lock, err := client.Lock(
		ctx,
		"demo:lock",
		dlock.WithTTL(5*time.Second),
		dlock.WithTryTimeout(3*time.Second),
		dlock.WithAutoRenew(true),
	)
	if err != nil {
		log.Fatal("acquire lock failed:", err)
	}
	defer lock.Unlock(ctx)

	fmt.Println("got lock:", lock.Key())

	// 演示“可重入”：同一 goroutine 再次获取同一把锁
	lock2, err := client.Lock(ctx, "demo:lock")
	if err != nil {
		log.Fatal("reenter lock failed:", err)
	}
	defer lock2.Unlock(ctx)

	fmt.Println("reentered lock")

	// 模拟一个比较长的任务，超过 TTL，但因为有自动续约不会丢锁
	time.Sleep(12 * time.Second)

	fmt.Println("business done")
}
