package review

import (
	"fmt"
	"time"
)

// ============================================================
// 补充模块：time 包
// ============================================================

func TimeDemo() {
	fmt.Println("=== time 包基础 ===")

	// 1. 当前时间
	now := time.Now()
	fmt.Printf("time.Now(): %v\n", now)

	// 2. Duration：时间间隔，本质是 int64 纳秒
	d := 3 * time.Second
	fmt.Printf("3s = %v = %d ns\n", d, d)

	// 3. Sleep
	fmt.Println("Sleep 100ms...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("醒来")

	// 4. Since：从某个时间点到现在过了多久
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("time.Since(start): %v\n", time.Since(start))
}

func TickerDemo() {
	fmt.Println("\n=== Ticker 定时触发 ===")

	// Ticker：每隔一段时间往 channel 发一个信号
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		fmt.Printf("  滴 %d (%v)\n", count, time.Now().Format("15:04:06"))
		if count >= 3 {
			break
		}
	}
	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C: // Ticker 的 channel
	//		fmt.Println("定时任务执行")
	//	case <-done:
	//		return
	//	}
	//}
}

func AfterDemo() {
	fmt.Println("\n=== time.After 超时控制 ===")

	// time.After：返回一个 channel，到时间了 channel 会收到信号
	// 配合 select 实现超时控制

	slowOp := func() string {
		time.Sleep(300 * time.Millisecond)
		return "done"
	}

	select {
	case result := <-time.After(200 * time.Millisecond):
		fmt.Printf("  结果: %v\n", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("  超时！")
		// 注意：上面这个 select 的逻辑有问题，应该用 goroutine
	}

	// 正确的用法：用 goroutine + channel
	fmt.Println("\n正确用法：goroutine + select：")
	ch := make(chan string, 1)
	go func() { ch <- slowOp() }()

	select {
	case result := <-ch:
		fmt.Printf("  结果: %v\n", result)
	case result := <-time.After(200 * time.Millisecond):
		fmt.Printf(" %v 超时！（慢操作被中断）", result)
	}
}

func DurationDemo() {
	fmt.Println("\n=== Duration 单位转换 ===")
	fmt.Printf("1s = %v\n", time.Second)
	fmt.Printf("1ms = %v\n", time.Millisecond)
	fmt.Printf("1min = %v\n", time.Minute)
}
