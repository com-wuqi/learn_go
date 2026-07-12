package practice

import "fmt"

// Check 验证练习 1-9（当前知识范围内的题目）
// 练习 10-14 需学完并发后手动验证
func Check() {
	fmt.Println("=== 练习 1-9 验证（当前可做）===")
	fmt.Println("=== 练习 10-14 待学完并发后自行验证 ===")
	fmt.Println()

	// 练习 1
	if r := Reverse("你好世界"); r == "界世好你" {
		fmt.Println("✅ 练习1  Reverse:     正确")
	} else {
		fmt.Printf("❌ 练习1  Reverse:     期望 '界世好你'，得到 %q\n", r)
	}

	// 练习 2
	if i, j := TwoSum([]int{2, 7, 11, 15}, 9); i == 0 && j == 1 {
		fmt.Println("✅ 练习2  TwoSum:      正确")
	} else {
		fmt.Printf("❌ 练习2  TwoSum:      期望 (0,1)，得到 (%d,%d)\n", i, j)
	}

	// 练习 3
	wc := WordCount("hello world hello go")
	if wc != nil && wc["hello"] == 2 && wc["world"] == 1 && wc["go"] == 1 {
		fmt.Println("✅ 练习3  WordCount:   正确")
	} else {
		fmt.Printf("❌ 练习3  WordCount:   期望 map[hello:2 world:1 go:1]，得到 %v\n", wc)
	}
	top := TopWords("a b a c b a d", 2)
	if len(top) == 2 && top[0] == "a" && top[1] == "b" {
		fmt.Println("✅ 练习3  TopWords:    正确")
	} else {
		fmt.Printf("❌ 练习3  TopWords:    期望 [a b]，得到 %v\n", top)
	}

	// 练习 4
	rb := NewRingBuffer(3)
	if rb != nil {
		_ = rb.Push(1)
		_ = rb.Push(2)
		_ = rb.Push(3)
		if rb.IsFull() {
			v, _ := rb.Pop()
			if v == 1 {
				fmt.Println("✅ 练习4  RingBuffer:  正确")
			} else {
				fmt.Printf("❌ 练习4  RingBuffer:  Pop 期望 1，得到 %d\n", v)
			}
		} else {
			fmt.Println("❌ 练习4  RingBuffer:  IsFull 返回 false")
		}
	} else {
		fmt.Println("❌ 练习4  RingBuffer:  NewRingBuffer 返回 nil")
	}

	// 练习 5
	fib := MemoFib()
	if r := fib(10); r == 55 {
		fmt.Println("✅ 练习5  MemoFib:     正确")
	} else {
		fmt.Printf("❌ 练习5  MemoFib:     fib(10) 期望 55，得到 %d\n", r)
	}

	// 练习 6
	r := MaxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	expected6 := []int{3, 3, 5, 5, 6, 7}
	same := len(r) == len(expected6)
	if same {
		for i := range r {
			if r[i] != expected6[i] {
				same = false
				break
			}
		}
	}
	if same {
		fmt.Println("✅ 练习6  MaxSliding:  正确")
	} else {
		fmt.Printf("❌ 练习6  MaxSliding:  期望 %v，得到 %v\n", expected6, r)
	}

	// 练习 7
	merged := MergeSorted([]int{1, 3, 5}, []int{2, 4, 6})
	if len(merged) == 6 && merged[0] == 1 && merged[5] == 6 {
		fmt.Println("✅ 练习7  MergeSorted: 正确")
	} else {
		fmt.Printf("❌ 练习7  MergeSorted: 期望 [1 2 3 4 5 6]，得到 %v\n", merged)
	}

	// 练习 8
	deduped := Dedup([]int{3, 1, 2, 1, 3, 4})
	if len(deduped) >= 4 && deduped[0] == 3 && deduped[1] == 1 {
		fmt.Println("✅ 练习8  Dedup:       正确")
	} else {
		fmt.Printf("❌ 练习8  Dedup:       期望 [3 1 2 4]，得到 %v\n", deduped)
	}

	// 练习 9
	flat := Flatten([][]int{{1, 2}, {3, 4}, {5}})
	if len(flat) == 5 && flat[0] == 1 && flat[4] == 5 {
		fmt.Println("✅ 练习9  Flatten:     正确")
	} else {
		fmt.Printf("❌ 练习9  Flatten:     期望 [1 2 3 4 5]，得到 %v\n", flat)
	}
}
