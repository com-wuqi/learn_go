package practice

import (
	"errors"
	"fmt"
)

// ============================================================
// 练习 1：反转字符串（支持中文）
// ============================================================

// [TODO] 实现 Reverse，输入 "你好世界" 输出 "界世好你"
// 提示：string 直接按字节反转会破坏多字节字符，需要用 []rune
func Reverse(s string) string {
	// TODO: 在这里实现
	return ""
}

// ============================================================
// 练习 2：两数之和
// ============================================================

// [TODO] 在 nums 中找两个数相加等于 target，返回它们的下标
// 如果不存在则返回 (-1, -1)，用 map 实现 O(n)
func TwoSum(nums []int, target int) (int, int) {
	// TODO: 在这里实现
	return -1, -1
}

// ============================================================
// 练习 3：词频统计与排序
// ============================================================

// WordCount 返回每个单词出现的次数
func WordCount(s string) map[string]int {
	// TODO: 在这里实现
	return nil
}

// TopWords 返回出现频率最高的 topN 个词，按频率降序排列
func TopWords(s string, topN int) []string {
	// TODO: 在这里实现
	return nil
}

// ============================================================
// 练习 4：环形缓冲区 (RingBuffer)
// ============================================================

// RingBuffer 是固定大小的环形队列
type RingBuffer struct {
	buf  []int
	head int // 读位置
	tail int // 写位置
	size int // 当前元素数
	cap  int // 容量
}

// [TODO] NewRingBuffer 创建一个容量为 capacity 的环形缓冲区
func NewRingBuffer(capacity int) *RingBuffer {
	// TODO
	return nil
}

// [TODO] Push 向缓冲区写入数据，如果满了返回 error
func (rb *RingBuffer) Push(val int) error {
	return errors.New("not implemented")
}

// [TODO] Pop 从缓冲区读取数据，如果为空返回 error
func (rb *RingBuffer) Pop() (int, error) {
	return 0, errors.New("not implemented")
}

// [TODO] IsEmpty 判断是否为空
func (rb *RingBuffer) IsEmpty() bool {
	return true
}

// [TODO] IsFull 判断是否已满
func (rb *RingBuffer) IsFull() bool {
	return false
}

// ============================================================
// 练习 5：并发安全计数器
// ============================================================

type SafeCounter struct {
	// TODO: 添加字段，用 sync.Mutex 保护
	val int
}

// [TODO] NewSafeCounter 创建计数器
func NewSafeCounter() *SafeCounter {
	return nil
}

// [TODO] Inc 加 1
func (c *SafeCounter) Inc() {
}

// [TODO] Value 读取当前值
func (c *SafeCounter) Value() int {
	return 0
}

// ============================================================
// 练习 6：超时重试函数
// ============================================================

// [TODO] RetryWithTimeout 在指定时间内重试 fn，直到成功或超时
// fn 返回 error 表示失败；如果超时返回超时错误
func RetryWithTimeout(maxRetries int, fn func() error) error {
	return errors.New("not implemented")
}

// ============================================================
// 练习 7：Worker Pool
// ============================================================

// [TODO] 实现一个固定数量的 worker pool
// - workers 个 goroutine 并发处理任务
// - 调用者通过 Submit 提交任务
// - 调用 Stop 后不再接收新任务，等待所有任务完成
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

// ============================================================
// 练习 8：简单 LRU Cache
// ============================================================

// [TODO] 实现 LRU 缓存（Least Recently Used）
// - 固定容量，超过容量时淘汰最久未使用的
// - 用 Map + 双向链表实现 O(1) Get/Put
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

// ============================================================
// 练习 9：交替打印
// ============================================================

// [TODO] 启动两个 goroutine，交替打印 1-10
// goroutine1 打印奇数，goroutine2 打印偶数，输出 "1 2 3 4 5 6 7 8 9 10"
// 提示：用两个 channel 或一个 channel 实现
func AlternatePrint() string {
	return ""
}

// ============================================================
// 练习 10：速查填空（基础知识）
// ============================================================

func Quiz() {
	fmt.Println("\n=== 速查填空：边看边回答 ===")

	// Q1: 以下代码输出什么？
	nums := []int{1, 2, 3}
	modify := func(s []int) {
		s[0] = 100
		s = append(s, 4) // 可能扩容，s 指向新数组
		s[1] = 200
	}
	modify(nums)
	fmt.Printf("Q1: nums=%v\n", nums) // [TODO] 你的答案: _______
	// 答: [100 2 3] — append 后 s 是新切片，原来的 nums 不变(除 s[0]=100)

	// Q2: defer 输出顺序？
	fmt.Print("Q2: ")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("%d ", i)
	}
	// [TODO] 你的答案: _______
	// 答: 2 1 0 (先进后出)

	// Q3: map 取的 key 不存在返回什么？
	m := map[string]int{"a": 1}
	fmt.Printf("\nQ3: m['b']=%d\n", m["b"])
	// [TODO] 你的答案: _______
	// 答: 0 (零值)

	// Q4: 以下是否会 panic？
	// var p *int
	// *p = 42
	// [TODO] 你的答案: _______
	// 答: 会 panic (nil pointer dereference)

	// Q5: cap 的作用？
	s := make([]int, 3, 10)
	fmt.Printf("Q5: len=%d cap=%d, append 3 个元素后是否扩容？\n", len(s), cap(s))
	s = append(s, 1, 2, 3)
	fmt.Printf("     len=%d cap=%d (cap 够大，未扩容)\n", len(s), cap(s))
	// cap 预分配避免多次扩容，提高性能

	// Q6: 如何安全检测类型断言？
	var x interface{} = "hello"
	// if s, ok := x.(string); ok { ... }
	// [TODO] 练习：用 comma-ok 模式检测 x 是否为 int
	if _, ok := x.(int); ok {
		fmt.Println("Q6: x 是 int")
	} else {
		fmt.Printf("Q6: x 不是 int，实际类型是 %T\n", x)
	}

	// Q7: 切片 a[1:3] 和 a[1:3:4] 的 cap 各是多少？
	a := []int{1, 2, 3, 4, 5}
	fmt.Printf("Q7: a[1:3] cap=%d, a[1:3:4] cap=%d\n", cap(a[1:3]), cap(a[1:3:4]))
	// [TODO] 答案: cap(a[1:3])=4 (5-1), cap(a[1:3:4])=3 (4-1)

	// Q8: named return 与 defer 交互
	q8 := func() (r int) {
		defer func() { r++ }()
		return 1
	}
	fmt.Printf("Q8: %d\n", q8())
	// [TODO] 答案: 2 (return 赋 r=1, defer 执行 r++, 最终返回 2)
}
