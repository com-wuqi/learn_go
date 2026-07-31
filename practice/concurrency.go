package practice

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================
// 并发练习（学完 goroutine/channel/sync 后做）
// ============================================================

// 练习 10：并发安全计数器
type SafeCounter struct {
	// TODO: 添加字段，用 sync.Mutex 保护
	val int
	mu  sync.Mutex
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{val: 0, mu: sync.Mutex{}}
}
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

// 练习 11：超时重试函数
func RetryWithTimeout(maxRetries int, fn func() error) error {
	return errors.New("not implemented")
}

// 练习 12：Worker Pool
type WorkerPool struct {
	wg     sync.WaitGroup
	jobs   chan func()
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
	workerPool := &WorkerPool{
		jobs: make(chan func(), workers),
		wg:   sync.WaitGroup{},
	}
	workerPool.ctx, workerPool.cancel = context.WithCancel(context.Background())
	for i := 0; i < workers; i++ {
		workerPool.wg.Add(1)
		go func(id int) {
			defer workerPool.wg.Done()
			fmt.Printf("worker %d starting\n", id)
			for {
				select {
				case <-workerPool.ctx.Done(): // 这里也是阻塞读
					return
				case job, ok := <-workerPool.jobs:
					// 这里是阻塞读，两个case竞争就绪，
					// default永远就绪，default里面阻塞读，无法感知外界Done
					if !ok {
						return
					}
					job()
				}
			}
		}(i)
	}
	return workerPool
}
func (wp *WorkerPool) Submit(fn func()) {
	wp.jobs <- fn
}
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.jobs)
}

// 练习 13：简单 LRU Cache
type LRUCache struct {
	// TODO: 添加字段
}

func NewLRUCache(capacity int) *LRUCache {
	return nil
}
func (c *LRUCache) Get(key string) (int, bool) {
	return 0, false
}
func (c *LRUCache) Put(key string, value int) {
}

// 练习 14：交替打印
func AlternatePrint() string {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch3 := make(chan int, 11)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			temp, ok := <-ch1
			if !ok {
				close(ch2)
				break
			}
			ch3 <- temp
			if temp >= 10 {
				close(ch2)
				break
			}
			ch2 <- temp + 1
		}
	}(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			temp, ok := <-ch2
			if !ok {
				close(ch1)
				break
			}
			ch3 <- temp
			if temp >= 10 {
				close(ch1)
				break
			}
			ch1 <- temp + 1
		}
	}(wg)
	wg.Wait()
	close(ch3)
	result := ""
	for ch := range ch3 {
		result += strconv.Itoa(ch) + " "
	}
	return strings.TrimSpace(result)
}

// ============================================================
// 练习 47：并发安全的 Map（用 sync.RWMutex）
// ============================================================
// [TODO] SafeMap 包装 map[string]string，Get 用读锁，Set/Delete 用写锁
// 提示: RLock/RUnlock for Get, Lock/Unlock for Set/Delete
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{data: make(map[string]string), mu: sync.RWMutex{}}
}

func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if val, ok := sm.data[key]; ok {
		return val, ok
	}
	return "", false
}

func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// ============================================================
// 练习 48：sync.Once — 延迟初始化
// ============================================================
// [TODO] GetConfig 首次调用时才从环境变量/文件加载配置，后续调用返回缓存的配置
// 提示: sync.Once.Do(func() { ... })
var (
	configOnce sync.Once
	cachedCfg  *ServerCfg
)

func GetConfig() *ServerCfg {
	configOnce.Do(func() {
		cachedCfg = &ServerCfg{}
	})
	return cachedCfg
}

// ============================================================
// 练习 49：Pipeline — 数据流水线
// ============================================================
// [TODO] SquarePipeline 输入整数切片，返回平方的切片
// stage1: gen(nums...) 将切片写入 channel
// stage2: square(in) 从 channel 读，计算平方，输出到 channel（复用 review/patterns.go 里的思路）
// 提示: 可以复用 review 里的 gen 和 square，自己写版本也行
func SquarePipeline(nums []int) []int {
	var result []int
	for num := range Square(Gen(nums)) {
		result = append(result, num)
	}
	return result
}
func Gen(nums []int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, num := range nums {
			ch <- num
		}
	}()
	return ch
}
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			out <- num * num
		}
	}()
	return out
}

// ============================================================
// 练习 50：Context 取消
// ============================================================
// [TODO] RunWithCancel 启动一个 goroutine 执行 fn，返回 ctx 和 cancel
// 调用方通过 cancel() 取消 goroutine，goroutine 通过 ctx.Done() 检测退出
// fn 内部应该是周期性工作 + select ctx.Done() 的模式
func RunWithCancel(fn func(context.Context)) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go fn(ctx)
	return ctx, cancel
}

// ============================================================
// 练习题 51：并行求平方和
// ============================================================
// [TODO] ParallelSum 启动 n 个 goroutine，每个计算一部分数字的平方，通过 channel 汇总
// 输入: nums=[1,2,3,4,5,6,7,8], goroutines=4
// 每个 goroutine 处理 2 个数，求平方，发到 channel，主 goroutine 累加
// 提示: 用 WaitGroup 等所有 worker 完成，close channel 后 range 累加
func ParallelSum(nums []int, goroutines int) int {
	out := make(chan int, goroutines)
	sumUp := 0
	var wg sync.WaitGroup
	slices := make([][]int, 0, goroutines)
	chunkSize := len(nums) / goroutines
	leftSize := len(nums) % goroutines // 剩余的
	for i := 0; i < len(nums); {
		if leftSize > 0 {
			slices = append(slices, nums[i:i+1+chunkSize])
			i += chunkSize + 1
			leftSize -= 1
		} else {
			slices = append(slices, nums[i:i+chunkSize])
			i += chunkSize
		}
	}
	for _, d := range slices {
		wg.Add(1)
		go func(d []int) {
			defer wg.Done()
			sum := 0
			for _, a := range d {
				sum += a * a
			}
			out <- sum
		}(d)
	}
	wg.Wait()
	close(out)
	for v := range out {
		sumUp += v
	}
	return sumUp
}

// ============================================================
// 练习题 52：带 Context 超时的 WorkerPool
// ============================================================
// [TODO] NewWorkerPoolWithTimeout 在 WorkerPool 基础上，整个池子有最大存活时间
// 超时后自动调用 Stop()，无需外部调用 cancel
func NewWorkerPoolWithTimeout(workers int, timeout time.Duration) *WorkerPool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	pool := NewWorkerPool(workers)
	go func() {
		<-ctx.Done()
		pool.Stop()
		cancel()
	}()
	return pool
}

// ============================================================
// 练习题 53：实现一个简易的并发限流器（Rate Limiter）
// ============================================================
// [TODO] RateLimiter 通过 buffered channel 限制并发数
// Acquire() 获取许可（阻塞直到有可用槽位）
// Release() 释放许可
// 提示: make(chan struct{}, limit) 作为信号量
type RateLimiter struct {
	// TODO
	bucket chan struct{}
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{bucket: make(chan struct{}, limit)}
}

func (rl *RateLimiter) Acquire() {
	rl.bucket <- struct{}{}
}

func (rl *RateLimiter) Release() {
	<-rl.bucket
}
