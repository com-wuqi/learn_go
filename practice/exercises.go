package practice

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// ============================================================
// 已掌握 ✅ — 基础数据结构与算法
// ============================================================

func Reverse(s string) string {
	var runes []rune
	for _, d := range s {
		runes = append(runes, d)
	}
	slices.Reverse(runes)
	return string(runes)
}

func TwoSum(nums []int, target int) (int, int) {
	hash := make(map[int]int)
	for i, d := range nums {
		if j, ok := hash[target-d]; ok {
			return j, i
		}
		hash[d] = i
	}
	return -1, -1
}

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
	s1 := strings.Fields(strings.TrimSpace(s))
	hash := make(map[string]int)
	for _, d := range s1 {
		if _, ok := hash[d]; ok {
			hash[d] = hash[d] + 1
		} else {
			hash[d] = 1
		}
	}
	sorted := make([]Pair, 0, len(hash))
	for k, v := range hash {
		sorted = append(sorted, Pair{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })
	var ans []string
	for i := 0; i < topN; i++ {
		ans = append(ans, sorted[i].str)
	}
	return ans
}

type RingBuffer struct {
	buf  []int
	head int
	tail int
	size int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{make([]int, capacity), 0, 0, 0}
}

func (rb *RingBuffer) Push(val int) error {
	if cap(rb.buf) == rb.size {
		return errors.New("buffer is full")
	}
	rb.buf[rb.tail] = val
	rb.tail = (rb.tail + 1) % cap(rb.buf)
	rb.size += 1
	return nil
}

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

func (rb *RingBuffer) IsEmpty() bool {
	return rb.size == 0
}

func (rb *RingBuffer) IsFull() bool {
	return rb.size == len(rb.buf)
}

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
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

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
// 速查填空（基础知识）
// ============================================================

func Quiz() {
	fmt.Println("\n=== 速查填空：边看边回答 ===")

	nums := []int{1, 2, 3}
	modify := func(s []int) {
		s[0] = 100
		s = append(s, 4)
		s[1] = 200
	}
	modify(nums)
	fmt.Printf("Q1: nums=%v\n", nums)
	// 答: [100 2 3]

	fmt.Print("Q2: ")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("%d ", i)
	}
	// 答: 2 1 0 (LIFO)

	m := map[string]int{"a": 1}
	fmt.Printf("\nQ3: m['b']=%d\n", m["b"])
	// 答: 0 (零值)

	// Q4: var p *int; *p = 42 → panic (nil pointer dereference)

	s := make([]int, 3, 10)
	fmt.Printf("Q5: len=%d cap=%d, append 3 个元素后是否扩容？\n", len(s), cap(s))
	s = append(s, 1, 2, 3)
	fmt.Printf("     len=%d cap=%d (cap 够大，未扩容)\n", len(s), cap(s))

	var x interface{} = "hello"
	if _, ok := x.(int); ok {
		fmt.Println("Q6: x 是 int")
	} else {
		fmt.Printf("Q6: x 不是 int，实际类型是 %T\n", x)
	}

	a := []int{1, 2, 3, 4, 5}
	fmt.Printf("Q7: a[1:3] cap=%d, a[1:3:4] cap=%d\n", cap(a[1:3]), cap(a[1:3:4]))
	// 答: cap(a[1:3])=4, cap(a[1:3:4])=3

	q8 := func() (r int) {
		defer func() { r++ }()
		return 1
	}
	fmt.Printf("Q8: %d\n", q8())
	// 答: 2
}
