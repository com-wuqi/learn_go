package review

import (
	"fmt"
	"slices"
	"sort"
)

// ============================================================
// 复习：数组、切片、Map
// ============================================================

// 编译时确定长度
var globalArr = [5]int{1, 2, 3, 4, 5}

func Collections() {
	fmt.Println("=== 复习：集合类型 ===")

	// --- 数组 ---
	var a1 [3]int              // [0,0,0]
	a2 := [3]int{10, 20, 30}   // 指定长度
	a3 := [...]int{1, 2, 3, 4} // 编译器推断长度
	var matrix [2][3]int       // 多维数组
	matrix[0] = [3]int{1, 2, 3}
	fmt.Printf("a1=%v len=%d\n", a1, len(a1))
	fmt.Printf("a2=%v\n", a2)
	fmt.Printf("a3=%v len=%d\n", a3, len(a3))
	fmt.Printf("matrix=%v\n", matrix)

	// 数组是值类型！赋值会复制整个数组
	a2copy := a2
	a2copy[0] = 999
	fmt.Printf("原数组a2[0]=%d, 复制后a2copy[0]=%d (值类型，互不影响)\n", a2[0], a2copy[0])

	// --- 切片 ---
	// slice 本质：ptr(指向底层数组) + len + cap
	// cap是容量，适用于已知大小的场景
	s1 := []int{1, 2, 3, 4, 5} // 字面量
	s2 := make([]int, 3, 5)    // make([]T, len, cap): [0,0,0] cap=5
	s3 := make([]int, 0, 4)    // len=0, cap=4

	fmt.Printf("\ns1=%v len=%d cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2=%v len=%d cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("s3=%v len=%d cap=%d\n", s3, len(s3), cap(s3))

	// append：当 len < cap 时不扩容；超过 cap 时分配新数组（通常是 2x）
	s3 = append(s3, 1, 2, 3)
	fmt.Printf("append 后: s3=%v len=%d cap=%d\n", s3, len(s3), cap(s3))
	s3 = append(s3, 4, 5) // 这里可能触发扩容
	fmt.Printf("扩容后: s3=%v len=%d cap=%d\n", s3, len(s3), cap(s3))

	// 切片操作：s[low:high:max] → len=high-low, cap=max-low
	original := []int{0, 1, 2, 3, 4, 5, 6, 7}
	sub := original[2:5]      // [2,3,4] len=3 cap=6（共享底层数组！）
	subCap := original[2:5:5] // [2,3,4] len=3 cap=3（限制 cap）
	fmt.Printf("\noriginal=%v\n", original)
	fmt.Printf("sub=original[2:5]     = %v len=%d cap=%d\n", sub, len(sub), cap(sub))
	fmt.Printf("subCap=original[2:5:5] = %v len=%d cap=%d\n", subCap, len(subCap), cap(subCap))
	// no max: 不关心共享
	// max = high: 隔离，append触发分配
	// max 够大 ：允许覆盖

	// 共享底层数组的陷阱：修改 sub 会影响 original
	sub[0] = 999
	fmt.Printf("修改 sub[0] 后: original=%v (也被改了!)\n", original)

	// 用三索引切片限制 cap 可以隔离
	original2 := []int{0, 1, 2, 3, 4, 5}
	isolated := original2[2:4:4] // cap=2，和 len 相同
	fmt.Printf("三索引隔离: isolated=%v len=%d cap=%d\n", isolated, len(isolated), cap(isolated))
	isolated = append(isolated, 100) // cap 不够，触发扩容，脱离原数组
	isolated[0] = 999
	fmt.Printf("扩容后修改: original2=%v (未受影响)\n", original2)

	// copy：min(len(src), len(dst)) 个元素
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	fmt.Printf("\ncopy: dst=%v (复制了 %d 个)\n", dst, n)

	// slices 标准库（Go 1.21+）
	fmt.Printf("slices.Reverse: ")
	slices.Reverse(src)
	fmt.Println(src)
	fmt.Printf("slices.Sort: ")
	slices.Sort(src)
	fmt.Println(src)

	// 删除元素（用 append + 切片操作）
	src = append(src[:2], src[3:]...) // 删除 index=2
	fmt.Printf("删除 index=2 后: %v\n", src)
	slices.Delete(src, 1, 2) // 不包含最后一个

	// --- Map ---
	fmt.Println("\n--- Map ---")
	m := map[string]int{
		"alice": 30,
		"bob":   25,
	}
	m["charlie"] = 35

	// 读取
	fmt.Printf("bob: %d\n", m["bob"])

	// comma-ok 检测 key 是否存在
	if age, ok := m["david"]; ok {
		fmt.Printf("david 存在，年龄=%d\n", age)
	} else {
		fmt.Println("david 不存在") // key 不存在返回零值 0
	}

	// 删除
	delete(m, "alice")
	fmt.Printf("删除 alice 后: %v len=%d\n", m, len(m))

	// map 遍历顺序是随机的！这是 Go 有意为之
	fmt.Print("遍历 map: ")
	for k, v := range m {
		fmt.Printf("%s:%d ", k, v)
	}
	fmt.Println(" (顺序随机)")

	// map 是非并发安全的
	// go func() { m["concurrent"] = 1 }() // 并发写会 panic!

	// 排序遍历 map：先取 key，排序，再按 key 取值
	fmt.Print("排序遍历 map：先取 key，排序，再按 key 取值:\n ")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s:%d ", k, m[k])
	}
	fmt.Println()

	// 嵌套结构：统计词频
	wordFreq := map[string]int{
		"hello": 3,
		"world": 2,
		"go":    5,
	}
	fmt.Printf("词频: %v\n", wordFreq)
	// 按频率排序
	type kv struct {
		k string
		v int
	}
	sorted := make([]kv, 0, len(wordFreq))
	for k, v := range wordFreq {
		sorted = append(sorted, kv{k, v})
	}
	forSortedFunc := func(i, j int) bool {
		return sorted[i].v > sorted[j].v // true: i 在 j 前面
	}
	sort.Slice(sorted, forSortedFunc)
	fmt.Print("按频率排序: ")
	for _, item := range sorted {
		fmt.Printf("%s:%d ", item.k, item.v)
	}
	fmt.Println()
}
