package review

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================
// 并发 Part 2：同步原语
// ============================================================

// --- Mutex：互斥锁，保护共享数据 ---

func MutexDemo() {
	fmt.Println("=== sync.Mutex ===")

	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			counter++ // 临界区：同一时间只有一个 goroutine 执行
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex 保护后: %d（1000 个 goroutine 各 +1）\n", counter)

	// ⚠️ 不加锁的话 counter 可能 < 1000（数据竞争）
}

// --- RWMutex：读写锁 ---

func RWMutexDemo() {
	fmt.Println("\n=== sync.RWMutex（读写锁）===")

	var rw sync.RWMutex
	data := map[string]int{"x": 0}
	var wg sync.WaitGroup

	// 读锁：多个 goroutine 可以同时持有
	// 写锁：独占，所有读锁和写锁都等它释放
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			rw.RLock() // 读锁
			fmt.Printf("  reader %d: x=%d\n", id, data["x"])
			rw.RUnlock()
			wg.Done()
		}(i)
	}

	wg.Add(1)
	go func() {
		rw.Lock() // 写锁：等所有读锁释放后才能写
		data["x"] = 42
		rw.Unlock()
		wg.Done()
	}()

	wg.Wait()
}

// --- Once：单次初始化 ---

type Singleton struct{ Name string }

var (
	instance *Singleton
	once     sync.Once
)

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{Name: "我是单例"}
		fmt.Println("  （初始化只执行一次）")
	})
	return instance
}

func OnceDemo() {
	fmt.Println("\n=== sync.Once（单例）===")

	for i := 0; i < 3; i++ {
		s := GetInstance()
		fmt.Printf("  获取: %s\n", s.Name)
	}
}

// --- atomic：原子操作 ---

func AtomicDemo() {
	fmt.Println("\n=== sync/atomic（原子操作）===")

	var counter int64 // atomic 要求 int64/int32
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&counter, 1) // 原子加一，不需要 Mutex
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Printf("atomic 计数: %d\n", counter)

	// CAS (Compare-And-Swap) — 乐观锁的基础
	var flag int32
	swapped := atomic.CompareAndSwapInt32(&flag, 0, 1) // flag==0 才改成 1
	fmt.Printf("CAS: flag=0→1 = %t\n", swapped)

	swapped = atomic.CompareAndSwapInt32(&flag, 0, 2) // flag 已经是 1，不匹配
	fmt.Printf("CAS: flag=0→2 = %t（失败）\n", swapped)
}

// --- 综合示例：并发安全的计数器 ---

type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func AtomicCounterDemo() {
	c := &AtomicCounter{}
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Printf("\n100k goroutine atomic 递增: %d（耗时 %v）\n", c.Value(), time.Since(start))
}
