package review

import (
	"context"
	"fmt"
	"time"
)

// ============================================================
// 并发 Part 3：context 包
// ============================================================
// context 解决了 goroutine 之间传递"取消信号"和"超时"的需求
// 核心方法: WithCancel, WithTimeout, WithValue

// --- 取消传播 ---

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done(): // context 被取消
			fmt.Printf("  worker %d: 收到取消，退出\n", id)
			return
		default:
			fmt.Printf("  worker %d: 工作中...\n", id)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func ContextCancelDemo() {
	fmt.Println("=== context.WithCancel（取消传播）===")

	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx, 1)
	go worker(ctx, 2)

	time.Sleep(1 * time.Second)
	fmt.Println("  主: 取消所有 worker")
	cancel() // 一键取消所有子 goroutine

	time.Sleep(200 * time.Millisecond) // 等 worker 退出
}

// --- 超时 ---

func ContextTimeoutDemo() {
	fmt.Println("\n=== context.WithTimeout（超时控制）===")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 好的习惯：用完后 cancel 释放资源

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("  操作完成")
	case <-ctx.Done():
		fmt.Printf("  超时: %v\n", ctx.Err()) // "context deadline exceeded"
	}
}

// --- 传值（trace ID / 请求元数据） ---

func ContextValueDemo() {
	fmt.Println("\n=== context.WithValue（传值）===")

	ctx := context.WithValue(context.Background(), "traceID", "abc-123")
	ctx = context.WithValue(ctx, "userID", "42")

	fmt.Printf("  traceID: %v\n", ctx.Value("traceID"))
	fmt.Printf("  userID: %v\n", ctx.Value("userID"))
	fmt.Printf("  不存在的key: %v\n", ctx.Value("nonexistent"))
}
