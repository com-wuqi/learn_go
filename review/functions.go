package review

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================
// 复习：函数、闭包、defer、错误处理
// ============================================================

// 多返回值和命名返回值
// Divide 除法运算，命名返回值可让函数签名更清晰
func Divide(a, b float64) (result float64, err error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	result = a / b
	return // 裸 return：返回命名返回值的当前值
}

// 变长参数
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 函数作为参数
func Apply(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, v := range nums {
		result[i] = fn(v)
	}
	return result
}

// 函数作为返回值：闭包工厂
func Counter() func() int {
	count := 0 // 逃逸到堆上
	return func() int {
		count++
		return count
	}
}

func Functions() {
	fmt.Println("=== 复习：函数、闭包、defer ===")

	// 多返回值
	r, err := Divide(10, 3)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("Divide(10,3) = %.2f\n", r)
	}

	// 检测除零
	if _, err := Divide(10, 0); err != nil {
		fmt.Printf("Divide(10,0) 错误: %v\n", err)
	}

	// 变长参数
	fmt.Printf("Sum(1,2,3,4,5) = %d\n", Sum(1, 2, 3, 4, 5))

	// 展开切片传入变长参数
	nums := []int{1, 2, 3}
	fmt.Printf("Sum(nums...) = %d\n", Sum(nums...))

	// 函数作为参数（map/filter/reduce 模式）
	doubled := Apply([]int{1, 2, 3, 4}, func(n int) int { return n * 2 })
	fmt.Printf("Apply double: %v\n", doubled)

	// 闭包：计数器
	fmt.Print("闭包计数器: ")
	c1 := Counter()
	c2 := Counter() // 每次调用 Counter() 创建独立的闭包
	for i := 0; i < 3; i++ {
		fmt.Printf("c1=%d ", c1())
	}
	fmt.Print("| ")
	for i := 0; i < 3; i++ {
		fmt.Printf("c2=%d ", c2())
	}
	fmt.Println("(c1 和 c2 独立)")

	// 闭包陷阱：循环变量捕获（Go 1.22 前经典 bug）
	fmt.Print("闭包陷阱演示: ")
	var funcs []func() int
	for i := 0; i < 3; i++ {
		i := i // 创建局部副本（Go 1.22+ 循环变量每次迭代都是新的，但这是好习惯）
		funcs = append(funcs, func() int { return i })
	}
	for _, fn := range funcs {
		fmt.Printf("%d ", fn())
	}
	fmt.Println("(Gp 1.22+ 也建议显式创建副本)")

	// defer：LIFO 栈顺序
	deferDemo()
}

func deferDemo() {
	fmt.Println("\n--- defer LIFO 演示 ---")

	// defer 立即求值参数，但延迟执行函数体
	name := "world"
	defer fmt.Printf("defer 1: %s (name 在 defer 声明时求值)\n", name)
	name = "golang"                         // 修改不影响 defer 1 的参数
	defer fmt.Printf("defer 2: %s\n", name) // 但 defer 2 声明时 name 已经是 "golang"

	// 闭包 defer：通过闭包引用的是最新值
	value := 0
	defer func() {
		fmt.Printf("defer 闭包: value=%d (引用的是最新值)\n", value)
	}()
	value = 42

	// defer 与 return 的执行顺序
	fmt.Printf("deferAndReturn 结果: %d (return 先求值，defer 后执行)\n", deferAndReturn())
}

// deferAndReturn 演示 defer 与命名返回值的交互
func deferAndReturn() (result int) {
	defer func() {
		result *= 2 // 修改命名返回值
		fmt.Printf("  defer 内: result=%d\n", result)
	}()
	return 5 // 先赋值 result=5，再执行 defer，最终返回 10
}

// ============================================================
// 复习：panic 和 recover
// ============================================================

func PanicRecover() {
	fmt.Println("\n=== 复习：panic 和 recover ===")

	// recover 必须在 defer 中调用才有效
	safeDivide := func(a, b int) (result int) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("捕获 panic: %v\n", r)
				result = -1
			}
		}()
		return a / b
	}

	fmt.Printf("safeDivide(10, 2) = %d\n", safeDivide(10, 2))
	fmt.Printf("safeDivide(10, 0) = %d (被 recover 救回来了)\n", safeDivide(10, 0))

	// 注意：不要滥用 panic，普通错误用 error 返回
}

// ============================================================
// 复习：泛型 (Go 1.18+)
// ============================================================

// 泛型约束
type Number interface {
	int | int64 | float64
}

// 泛型函数
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 泛型 Map 函数
func Map[T any, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, v := range items {
		result[i] = fn(v)
	}
	return result
}

func GenericsReview() {
	fmt.Println("\n=== 复习：泛型 ===")

	fmt.Printf("Max(3, 7) = %d\n", Max(3, 7))
	fmt.Printf("Max(3.14, 2.71) = %.2f\n", Max(3.14, 2.71))

	nums := []int{1, 2, 3, 4}
	strs := Map(nums, func(n int) string {
		return fmt.Sprintf("id_%d", n)
	})
	fmt.Printf("Map int→string: %v\n", strs)
}

// ============================================================
// 复习：错误处理模式
// ============================================================

// 自定义错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("验证失败: 字段=%s, 原因=%s", e.Field, e.Message)
}

// 错误包装 (Go 1.13+)
func Validate(account string, balance float64) error {
	if account == "" {
		return &ValidationError{Field: "account", Message: "不能为空"}
	}
	if balance < 0 {
		return fmt.Errorf("余额验证失败: %w", &ValidationError{Field: "balance", Message: "不能为负数"})
	}
	return nil
}

func ErrorsReview() {
	fmt.Println("\n=== 复习：错误处理 ===")

	// errors.Is：判断错误链中是否包含特定错误
	err := Validate("", 100)
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("类型断言成功: %v\n", valErr)
	}

	// errors.As：提取特定类型的错误
	err2 := Validate("user1", -50)
	if err2 != nil {
		fmt.Printf("错误: %v\n", err2)
		var vErr *ValidationError
		if errors.As(err2, &vErr) {
			fmt.Printf("提取到 ValidationError: field=%s\n", vErr.Field)
		}
	}

	// 常见错误模式
	fmt.Println("常见错误: time.Sleep 的 duration 必须是正数")
	_ = time.Now()
}
