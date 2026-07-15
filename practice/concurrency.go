package practice

import "errors"

// ============================================================
// 并发练习（学完 goroutine/channel/sync 后做）
// ============================================================

// 练习 10：并发安全计数器
type SafeCounter struct {
	// TODO: 添加字段，用 sync.Mutex 保护
	val int
}

func NewSafeCounter() *SafeCounter {
	return nil
}
func (c *SafeCounter) Inc() {
}
func (c *SafeCounter) Value() int {
	return 0
}

// 练习 11：超时重试函数
func RetryWithTimeout(maxRetries int, fn func() error) error {
	return errors.New("not implemented")
}

// 练习 12：Worker Pool
type WorkerPool struct {
	// TODO: 添加字段
}

func NewWorkerPool(workers int) *WorkerPool {
	return nil
}
func (wp *WorkerPool) Submit(fn func()) {
}
func (wp *WorkerPool) Stop() {
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
	return ""
}
