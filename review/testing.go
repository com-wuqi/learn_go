package review

import "fmt"

// ============================================================
// 模块：panic/recover + 测试基础
// ============================================================

// --- panic & recover ---
// panic = 程序遇到无法继续的错误，立即停止当前函数，沿调用栈向上"炸"
// recover = 在 defer 中"接住" panic，程序继续运行
// 原则：不要用 panic 做常规错误处理，只在"真正不可恢复"时用

func PanicDemo() {
	fmt.Println("=== panic & recover ===")

	// recover 必须在 defer 里调用
	safeDivide := func(a, b int) (result int) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  recover 接住了: %v\n", r)
				result = -1 // 返回默认值
			}
		}()
		return a / b // b=0 时 panic
	}

	fmt.Printf("  10/2 = %d\n", safeDivide(10, 2))
	fmt.Printf("  10/0 = %d (被 recover 救回)\n", safeDivide(10, 0))

	// 对比：常规错误用 error 返回，不要用 panic
	// ✅ func divide(a, b int) (int, error)
	// ❌ func divide(a, b int) int { if b==0 { panic("zero") } }
}

// --- 测试基本结构 ---
// Go 测试放到 *_test.go 文件中，函数名 TestXxx(t *testing.T)
// 运行: go test ./... 或 go test -v ./practice/

// 表驱动测试是 Go 社区标准写法：
// tests := []struct{ name, input, expected }{ ... }
// for _, tt := range tests {
//     t.Run(tt.name, func(t *testing.T) {
//         result := MyFunc(tt.input)
//         if result != tt.expected { t.Errorf(...) }
//     })
// }

// --- 基准测试 ---
// 函数名 BenchmarkXxx(b *testing.B)
// for i := 0; i < b.N; i++ { ... }  // b.N 自动调整
// 运行: go test -bench=. ./practice/

func TestExplain() {
	fmt.Println("\n`go test` 怎么用？")
	fmt.Println("  go test ./...              # 跑所有包的测试")
	fmt.Println("  go test -v ./practice/     # 详细输出")
	fmt.Println("  go test -run TestKVStore   # 只跑匹配的测试")
	fmt.Println("  go test -bench=.           # 跑基准测试")
	fmt.Println("  go test -cover ./...       # 查看代码覆盖率")
}
