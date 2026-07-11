package review

import (
	"fmt"
	"math/rand"
	"time"
)

// ============================================================
// 复习：if-else、switch、for、label
// ============================================================

func FlowControl() {
	fmt.Println("=== 复习：控制流 ===")

	// if 带初始化语句（常见于错误处理）
	if score := 85; score >= 90 {
		fmt.Println("A")
	} else if score >= 80 {
		fmt.Println("B") // 输出 B
	} else {
		fmt.Println("C")
	}
	// score 作用域仅限于 if-else 块

	// switch：不用写 break，不会 fall-through
	switch day := time.Now().Weekday(); day {
	case time.Saturday, time.Sunday: // 多值
		fmt.Println("周末！")
	case time.Monday:
		fmt.Println("周一...")
		fallthrough // 显式穿透到下一个 case
	default:
		fmt.Println("工作日")
	}

	// 无表达式的 switch（相当于 if-else 链）
	t := time.Now().Hour()
	switch {
	case t < 12:
		fmt.Println("上午")
	case t < 18:
		fmt.Println("下午")
	default:
		fmt.Println("晚上")
	}

	// type switch
	printType := func(v interface{}) {
		switch v := v.(type) {
		case int:
			fmt.Printf("int: %d\n", v)
		case string:
			fmt.Printf("string: %q (长度=%d)\n", v, len(v))
		case bool:
			fmt.Printf("bool: %t\n", v)
		case []int:
			fmt.Printf("[]int: %v\n", v)
		default:
			fmt.Printf("未知类型: %T\n", v)
		}
	}
	printType(42)
	printType("hello")
	printType(true)
	printType([]int{1, 2, 3})
	printType(3.14)

	// for 循环的 4 种写法
	fmt.Println("\n--- for 循环 ---")
	// 1. 经典三段式
	for i := 0; i < 3; i++ {
		fmt.Printf("三段式: %d ", i)
	}
	fmt.Println()

	// 2. while 式
	j := 3
	for j > 0 {
		fmt.Printf("while: %d ", j)
		j--
	}
	fmt.Println()

	// 3. 无限循环
	// for { ... break }

	// 4. for range
	fmt.Print("range slice: ")
	for i, v := range []string{"a", "b", "c"} {
		fmt.Printf("[%d]=%s ", i, v)
	}
	fmt.Println()

	fmt.Print("range map: ")
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	for k, v := range m {
		fmt.Printf("%s:%d ", k, v)
	}
	fmt.Println("(注意：map 遍历顺序随机!)")

	// range string：按 rune 遍历
	fmt.Printf("range 中文字符串: ")
	for i, r := range "你好Go" {
		fmt.Printf("[byte%d]=%c ", i, r)
	}
	fmt.Println()

	// break 和 continue
	fmt.Println("\n--- break/continue ---")
	for i := 0; i < 5; i++ {
		if i == 2 {
			continue // 跳过本次
		}
		if i == 4 {
			break // 跳出循环
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// label：控制多层循环
	fmt.Println("\n--- label ---")
	fmt.Println("九九乘法表（跳过对角线和超过 30 的项）:")
OUTER:
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			if i*j > 30 {
				break OUTER // 跳出外层循环
			}
			if i == j {
				continue // 跳过对角线
			}
			fmt.Printf("%dx%d=%-2d ", i, j, i*j)
		}
		fmt.Println()
	}
}

// ============================================================
// 复习：随机数与 make
// ============================================================

func RandomDemo() {
	fmt.Println("\n=== 复习：随机数 ===")

	// Go 1.20+ 不需要手动 seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Print("5 个随机数 [0,100): ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", rng.Intn(100))
	}
	fmt.Println()

	// 随机排列
	nums := []int{1, 2, 3, 4, 5}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	fmt.Printf("随机排列: %v\n", nums)
}
