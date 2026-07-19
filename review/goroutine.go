package review

import (
	"fmt"
	"time"
)

// ============================================================
// 并发 Part 1：goroutine + channel + select
// ============================================================
// goroutine: 轻量级用户态线程，一个 Go 程序可以轻松跑上万个
// channel:   goroutine 之间的通信管道，"不要通过共享内存来通信，通过通信来共享内存"
// select:    同时监听多个 channel，哪个先有数据就处理哪个

// --- 示例 1：启动 goroutine ---

func GoroutineDemo() {
	fmt.Println("=== goroutine 基础 ===")

	go func() {
		fmt.Println("  我在另一个 goroutine 里")
	}()

	time.Sleep(50 * time.Millisecond) // 等 goroutine 跑完（实际代码用 WaitGroup）
	fmt.Println("  主 goroutine 结束")
}

// --- 示例 2：无缓冲 channel（同步） ---

func UnbufferedChanDemo() {
	fmt.Println("\n=== 无缓冲 channel（同步传递）===")

	ch := make(chan string) // 无缓冲：发送方必须等接收方准备好

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "hello" // 阻塞直到有人接收
		fmt.Println("  发送完成")
	}()

	fmt.Println("  等待接收...")
	msg := <-ch // 阻塞直到有人发送
	fmt.Printf("  收到: %s\n", msg)
}

// --- 示例 3：有缓冲 channel（异步，容量 2） ---

func BufferedChanDemo() {
	fmt.Println("\n=== 有缓冲 channel（异步）===")

	ch := make(chan int, 2) // 缓冲 2：可以塞 2 个值而不阻塞

	ch <- 1 // 不阻塞
	ch <- 2 // 不阻塞
	// ch <- 3 // 阻塞！缓冲区满了

	fmt.Printf("  取出: %d\n", <-ch)
	fmt.Printf("  取出: %d\n", <-ch)
}

// --- 示例 4：channel 关闭 + range ---

func CloseRangeDemo() {
	fmt.Println("\n=== close + range ===")

	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch) // 关闭后不能再发送，但可以继续接收

	// range 会在 channel 关闭且缓冲读完后自动退出
	for v := range ch {
		fmt.Printf("  %d ", v)
	}
	fmt.Println()

	// 检测 channel 是否关闭
	v, ok := <-ch
	fmt.Printf("  关闭后读取: v=%d, ok=%t\n", v, ok) // 0, false
}

// --- 示例 5：select 多路复用 ---

func SelectDemo() {
	fmt.Println("\n=== select 多路复用 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "慢消息"
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "快消息"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  ch1: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("  ch2: %s\n", msg)
		}
	}
}

// --- 示例 6：select 的 default（非阻塞） ---

func SelectDefaultDemo() {
	fmt.Println("\n=== select default（非阻塞尝试）===")

	ch := make(chan int, 1)

	select {
	case ch <- 42:
		fmt.Println("  发送成功（buffer 有空位）")
	default:
		fmt.Println("  发送失败（buffer 满）")
	}

	// 非阻塞接收
	select {
	case v := <-ch:
		fmt.Printf("  接收成功: %d\n", v) // 输出这个
	default:
		fmt.Println("  接收失败（无数据）")
	}
}

// --- 示例 7：生产者-消费者 ---

func ProducerConsumerDemo() {
	fmt.Println("\n=== 生产者-消费者 ===")

	jobs := make(chan int, 5)
	done := make(chan bool)

	// 消费者
	go func() {
		for job := range jobs { // range 在 close 后自动结束
			fmt.Printf("  处理: %d\n", job)
			time.Sleep(100 * time.Millisecond)
		}
		done <- true
	}()

	// 生产者
	for i := 1; i <= 3; i++ {
		jobs <- i
		fmt.Printf("  提交: %d\n", i)
	}
	close(jobs) // 通知消费者不再有任务

	<-done // 等消费者完成
	fmt.Println("  全部完成")
}
