package practice

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// ============================================================
// 第一部分：当前可做（已掌握知识点）
// ============================================================

// ============================================================
// 练习 1：反转字符串（支持中文）
// [TODO] 实现 Reverse，输入 "你好世界" 输出 "界世好你"
// 提示：string 直接按字节反转会破坏多字节字符，需要用 []rune
// ============================================================
func Reverse(s string) string {
	var runes []rune
	for _, d := range s {
		runes = append(runes, d)
	}
	slices.Reverse(runes)
	return string(runes)
}

// ============================================================
// 练习 2：两数之和
// ============================================================

// [TODO] 在 nums 中找两个数相加等于 target，返回它们的下标
// 如果不存在则返回 (-1, -1)，用 map 实现 O(n)
func TwoSum(nums []int, target int) (int, int) {
	hash := make(map[int]int)
	for i, d := range nums {
		if j, ok := hash[target-d]; ok {
			return j, i
		}
		hash[d] = i
		// 这里需要先查后存，防止访问自身
	}
	return -1, -1
}

// ============================================================
// 练习 3：词频统计与排序
// WordCount 返回每个单词出现的次数
// TopWords 返回出现频率最高的 topN 个词，按频率降序排列
// ============================================================
func WordCount(s string) (hash map[string]int) {
	s1 := strings.Fields(strings.TrimSpace(s))
	hash = make(map[string]int)
	for _, d := range s1 {
		if _, ok := hash[d]; ok {
			hash[d] = hash[d] + 1
		} else {
			hash[d] = 1
		}
	}
	return
}

type Pair struct {
	str   string
	count int
}

func TopWords(s string, topN int) []string {
	// string -> map
	s1 := strings.Fields(strings.TrimSpace(s))
	hash := make(map[string]int)
	for _, d := range s1 {
		if _, ok := hash[d]; ok {
			hash[d] = hash[d] + 1
		} else {
			hash[d] = 1
		}
	}
	// msp -> []KV(sorted)
	sorted := make([]Pair, 0, len(hash))
	for k, v := range hash {
		sorted = append(sorted, Pair{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })
	// sorted []KV -> []string
	var ans []string
	for i := 0; i < topN; i++ {
		ans = append(ans, sorted[i].str)
	}
	return ans
}

// ============================================================
// 练习 4：环形缓冲区 (RingBuffer)
// RingBuffer 是固定大小的环形队列（消息队列基础）
// ============================================================

type RingBuffer struct {
	buf  []int
	head int // 读位置
	tail int // 写位置
	size int // 当前大小
}

// [TODO] NewRingBuffer 创建一个容量为 capacity 的环形缓冲区
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{make([]int, capacity), 0, 0, 0}
}

// [TODO] Push 向缓冲区写入数据，如果满了返回 error
func (rb *RingBuffer) Push(val int) error {
	if cap(rb.buf) == rb.size {
		return errors.New("buffer is full")
	}
	rb.buf[rb.tail] = val
	rb.tail = (rb.tail + 1) % cap(rb.buf)
	rb.size += 1
	return nil

}

// [TODO] Pop 从缓冲区读取数据，如果为空返回 error
func (rb *RingBuffer) Pop() (int, error) {
	if rb.size == 0 {
		return -1, errors.New("buffer is empty")
	}
	var val int
	val = rb.buf[rb.head]
	rb.head = (rb.head + 1) % cap(rb.buf)
	rb.size -= 1
	return val, nil
}

// [TODO] IsEmpty 判断是否为空
func (rb *RingBuffer) IsEmpty() bool {
	if rb.size == 0 {
		return true
	}
	return false
}

// [TODO] IsFull 判断是否已满
func (rb *RingBuffer) IsFull() bool {
	if rb.size == len(rb.buf) {
		return true
	}
	return false
}

// ============================================================
// 练习 5：斐波那契带记忆化
// 提示：闭包内维护一个 map 缓存已计算的值，fib(0)=0, fib(1)=1
// 每次计算前先查缓存 map，有则跳过 => 分布式缓存的微型版本
// ============================================================
func MemoFib() func(n int) int {
	cache := make(map[int]int)
	cache[0] = 0
	cache[1] = 1
	var fib func(n int) int
	fib = func(n int) int {
		if d, ok := cache[n]; ok {
			return d
		}
		cache[n] = fib(n-1) + fib(n-2)
		return cache[n]
	}
	return fib
}

// ============================================================
// 练习 6：滑块窗口最大值
// 输入 []int{1,3,-1,-3,5,3,6,7}, k=3
// 窗口每次右移一步，记录窗口内最大值：
//
//	[1,3,-1] → 3, [3,-1,-3] → 3, [-1,-3,5] → 5, [-3,5,3] → 5,
//	[5,3,6] → 6, [3,6,7] → 7
//
// 输出: [3,3,5,5,6,7]
// ============================================================
func MaxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 {
		return nil
	}
	result := make([]int, 0, len(nums)-k+1)

	for i := 0; i <= len(nums)-k; i++ {
		// 找出 window [i, i+k) 里的最大值
		maxVal := nums[i]
		for j := i + 1; j < i+k; j++ {
			if nums[j] > maxVal {
				maxVal = nums[j]
			}
		}
		result = append(result, maxVal)
	}
	return result
}

// ============================================================
// 练习 7：合并两个有序切片
// 输入: [1,3,5], [2,4,6] → 输出: [1,2,3,4,5,6]
// 提示：两个指针，谁小谁先进结果（归并排序的 merge 步骤）
// ============================================================
func MergeSorted(a, b []int) []int {
	// TODO: 在这里实现
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...) // 维护末尾
	result = append(result, b[j:]...)
	return result
}

// ============================================================
// 练习 8：切片去重（保持原顺序）
// 输入: [3,1,2,1,3,4] → 输出: [3,1,2,4]
// 提示：用 map 记录已见过的元素，遍历原切片，没见过的才 append
// ============================================================
func Dedup(s []int) []int {
	hash := make(map[int]bool)
	result := make([]int, 0, len(s))
	for _, d := range s {
		if !hash[d] {
			hash[d] = true
			result = append(result, d)
		}
	}
	return result
}

// ============================================================
// 练习 9：展平嵌套切片
// 输入: [][]int{{1,2}, {3,4}, {5}} → 输出: []int{1,2,3,4,5}
// ============================================================
func Flatten(nested [][]int) []int {
	var result []int
	for _, n := range nested {
		for _, d := range n {
			result = append(result, d)
		}
	}
	return result
}

// ============================================================
// 第二部分：学完并发后再做（goroutine / channel / sync）
// ============================================================

// ============================================================
// 练习 10：并发安全计数器
// [TODO] 添加字段，用 sync.Mutex 保护
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
// 练习 11：超时重试函数
// [TODO] RetryWithTimeout 在指定时间内重试 fn，直到成功或超时
// fn 返回 error 表示失败；如果超时返回超时错误
// 提示：time.Ticker 或 time.After；也可简单限定重试次数
// ============================================================
func RetryWithTimeout(maxRetries int, fn func() error) error {
	return errors.New("not implemented")
}

// ============================================================
// 练习 12：Worker Pool
// [TODO] 实现一个固定数量的 worker pool
// - workers 个 goroutine 并发处理任务
// - 调用者通过 Submit 提交任务
// - 调用 Stop 后不再接收新任务，等待所有任务完成
// 提示：用 channel 发任务，sync.WaitGroup 等完成
// ============================================================
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
// 练习 13：简单 LRU Cache
// [TODO] 实现 LRU 缓存（Least Recently Used）
// - 固定容量，超过容量时淘汰最久未使用的
// - 用 Map + 双向链表实现 O(1) Get/Put
// ============================================================
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
// 练习 14：交替打印
// [TODO] 启动两个 goroutine，交替打印 1-10
// goroutine1 打印奇数，goroutine2 打印偶数，输出 "1 2 3 4 5 6 7 8 9 10"
// 提示：用两个 channel 或一个 channel 实现
// ============================================================
func AlternatePrint() string {
	return ""
}

// ============================================================
// 第三部分：接口练习（学完 review/interfaces.go 后做）
// ============================================================

// ============================================================
// 练习 15：实现 sort.Interface 对字符串按长度排序
// 输入: []string{"a", "abc", "ab"} → 输出: []string{"a", "ab", "abc"}
// 提示：定义 type ByLen []string，实现 Len/Less/Swap 三个方法
// ============================================================
type ByLen []string

func (b ByLen) Len() int           { return 0 }
func (b ByLen) Less(i, j int) bool { return false }
func (b ByLen) Swap(i, j int)      {}

// ============================================================
// 练习 16：实现一个 KVStore 接口，分别用内存 map 和文件实现
// ============================================================
type KVStore interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

// [TODO] 内存实现
type MemStore struct {
	// TODO: 添加字段
}

func NewMemStore() *MemStore {
	return nil
}
func (m *MemStore) Get(key string) (string, error)     { return "", nil }
func (m *MemStore) Set(key string, value string) error { return nil }
func (m *MemStore) Delete(key string) error            { return nil }

// [TODO] 文件实现（用 JSON 持久化到文件）
type FileStore struct {
	// TODO: 添加字段（文件路径、内存缓存 map）
}

func NewFileStore(path string) *FileStore {
	return nil
}
func (f *FileStore) Get(key string) (string, error)     { return "", nil }
func (f *FileStore) Set(key string, value string) error { return nil }
func (f *FileStore) Delete(key string) error            { return nil }

// ============================================================
// 练习 17：类型断言 — 用 comma-ok 安全地取出 Animal 的具体类型
// ============================================================

type AnimalI interface {
	Speak() string
}

type DogI struct{ Name string }
type CatI struct{ Name string }

func (d DogI) Speak() string { return d.Name + ": 汪汪" }
func (c CatI) Speak() string { return c.Name + ": 喵喵" }

// 输入: []AnimalI{DogI{"旺财"}, CatI{"咪咪"}, DogI{"大黄"}}
// 输出: []string{"旺财: 汪汪", "咪咪: 喵喵", "大黄: 汪汪"}
func DescribeAnimals(animals []AnimalI) []string {
	// TODO: 遍历 animals，类型断言判断是 Dog 还是 Cat
	return nil
}

// ============================================================
// 练习 18：用接口实现一个简单的 Plugin 模式
// Plugin 接口定义 Process() 方法，两个插件实现：UppercasePlugin（转大写）、
// ReversePlugin（反转），Pipeline 接收一串插件依次执行
// ============================================================
type Plugin interface {
	Process(s string) string
}

type UppercasePlugin struct{}

func (u UppercasePlugin) Process(s string) string { return "" }

type ReversePlugin struct{}

func (r ReversePlugin) Process(s string) string { return "" }

// [TODO] Pipeline 依次执行 plugins 中的每个插件
func RunPipeline(input string, plugins []Plugin) string {
	return ""
}

// ============================================================
// 速查填空（基础知识）
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
