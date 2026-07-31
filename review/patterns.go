package review

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 并发 Part 3：Pipeline、Fan-out/Fan-in
// ============================================================

// --- Pipeline：多个 stage 串联，每个 stage 是一个 goroutine ---
// 数据像流水线一样经过 gen → square → print

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func PipelineDemo() {
	fmt.Println("=== Pipeline（流水线）===")

	// gen → square → print
	for n := range square(gen(2, 3, 4)) {
		fmt.Printf("  %d ", n) // 4 9 16
	}
	fmt.Println()
}

// --- Fan-out：一个 channel 分给多个 worker ---
// --- Fan-in：多个 channel 合并成一个 ---

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs)) // 计数
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FanOutFanInDemo() {
	fmt.Println("\n=== Fan-out / Fan-in ===")

	input := gen(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out: 把同一个 input 分给 3 个 square worker
	w1 := square(input)
	w2 := square(input)
	w3 := square(input)

	// Fan-in: 3 个 worker 的结果合并到一个 channel
	for n := range merge(w1, w2, w3) {
		fmt.Printf("%d ", n)
	}
	fmt.Println()
}

// --- Worker Pool：固定数量 goroutine 处理任务队列 ---

func WorkerPoolDemo() {
	fmt.Println("\n=== Worker Pool ===")

	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// 启动 3 个 worker
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("  worker %d: 处理 %d\n", id, job)
				time.Sleep(100 * time.Millisecond)
			}
		}(w)
	}

	// 投递 5 个任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // 通知 worker 没有更多任务
	wg.Wait()
	fmt.Println("  全部完成")
}
