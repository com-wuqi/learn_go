package review

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ============================================================
// 复习 1：变量、常量、基本类型、字符串
// ============================================================

// const 复习：类型化 vs 非类型化常量
const (
	typedInt       int = 42
	untypedFloat       = 3.14 // 非类型化，可以赋值给任何浮点类型
	untypedComplex     = 1 + 2i
	maxUint8           = 1<<8 - 1 // 表达式常量
)

// iota 枚举：分布式系统中状态码、协议版本常用
const (
	StatusPending  = iota // 0
	StatusRunning         // 1
	StatusSuccess         // 2
	StatusFailed          // 3
	StatusCanceled        // 4
)

const (
	_  = iota             // 跳过 0
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20
	GB                    // 1 << 30
)

// type 类型定义 vs 类型别名
type UserID int64    // 新类型，和 int64 不能直接运算
type AliasID = int64 // 别名，完全等价于 int64

func Basics() {
	fmt.Println("=== 复习：变量、常量、类型 ===")

	// 变量声明的 4 种方式
	var a int  // 声明，零值 0
	var b = 10 // 类型推断
	var c = 20 // 显式类型
	d := 30    // 短声明（仅函数内可用）
	// 多重赋值
	a, b = b, a // 交换：a=10, b=0
	fmt.Println(a, b, c, d)

	// 零值表
	var (
		zeroInt     int     // 0
		zeroFloat   float64 // 0.0
		zeroBool    bool    // false
		zeroString  string  // ""
		zeroPointer *int    // nil
		zeroSlice   []int   // nil
	)
	fmt.Printf("零值: int=%d float=%f bool=%t str=%q ptr=%v slice=%v\n",
		zeroInt, zeroFloat, zeroBool, zeroString, zeroPointer, zeroSlice)

	// 指针：&取地址，*解引用
	x := 42
	px := &x
	*px = 100 // 通过指针修改 x
	fmt.Printf("x=%d, *px=%d\n", x, *px)

	// 类型转换（Go 必须显式转换）
	var f = float64(x)
	var u = uint(f)
	fmt.Printf("float64=%f, uint=%d\n", f, u)

	// UserID 是新类型，需要显式转换
	var uid UserID = 100
	var raw = int64(uid) + 1
	_ = raw
	fmt.Printf("uid=%d (type: %T)\n", uid, uid)

	// 别名不需要转换
	var aid AliasID = 200
	var raw2 = aid + 1 // 直接运算
	fmt.Printf("aid=%d, raw2=%d\n", aid, raw2)

	// 字符串操作
	s := "  Hello, Go 语言!  "
	fmt.Printf("原始: %q\n", s)
	fmt.Printf("HasPrefix 'Hello': %t\n", strings.HasPrefix(strings.TrimSpace(s), "Hello"))
	fmt.Printf("Contains 'Go': %t\n", strings.Contains(s, "Go"))
	fmt.Printf("Index '语': %d\n", strings.Index(s, "语"))               // 返回字节索引！注意中文占 3 字节
	fmt.Printf("Replace: %s\n", strings.Replace(s, "Go", "Golang", 1)) // n 是个数
	fmt.Printf("Split: %v\n", strings.Split(strings.TrimSpace(s), " "))
	fmt.Printf("Fields: %v\n", strings.Fields(s))
	fmt.Printf("ToUpper: %s\n", strings.ToUpper(s))

	// Unicode 陷阱：len 返回字节数，不是字符数
	cn := "你好世界"
	fmt.Printf("len(%q) = %d (字节数), rune 数 = %d\n", cn, len(cn), len([]rune(cn)))

	// strconv 数字转换（需要导入 strconv）
	//n, err := strconv.Atoi("123")   // string → int
	//s2 := strconv.Itoa(456)         // int → string
}

// ============================================================
// 复习 2：位运算和数值类型
// ============================================================

func BitOps() {
	fmt.Println("\n=== 复习：位运算 ===")

	a := uint8(0b00001111) // 15
	b := uint8(0b00110011) // 51

	fmt.Printf("a = %08b (%d)\n", a, a)
	fmt.Printf("b = %08b (%d)\n", b, b)
	fmt.Printf("a & b  (AND)      = %08b (%d)\n", a&b, a&b)
	fmt.Printf("a | b  (OR)       = %08b (%d)\n", a|b, a|b)
	fmt.Printf("a ^ b  (XOR)      = %08b (%d)\n", a^b, a^b)
	fmt.Printf("^a     (NOT)      = %08b (%d)\n", ^a, ^a)
	fmt.Printf("a << 2 (左移)     = %08b (%d)\n", a<<2, a<<2)
	fmt.Printf("b >> 2 (右移)     = %08b (%d)\n", b>>2, b>>2)
	fmt.Printf("a &^ b (AND NOT/清0)= %08b (%d)\n", a&^b, a&^b) // 将 b 中为 1 的位在 a 中清零
}

// ============================================================
// 复习 3：浮点数和数学运算
// ============================================================

// ============================================================
// 速查：Printf 常用占位符
// ============================================================

/*
  %v	默认格式                    fmt.Printf("%v", 42)          → 42
  %+v	结构体含字段名               fmt.Printf("%+v", point)     → {X:1 Y:2}
  %#v	Go 语法表示                  fmt.Printf("%#v", point)     → image.Point{X:1, Y:2}
  %T	类型                        fmt.Printf("%T", 42)          → int

  %d	十进制整数                   fmt.Printf("%d", 42)          → 42
  %b	二进制                      fmt.Printf("%b", 42)          → 101010
  %o	八进制                      fmt.Printf("%o", 42)          → 52
  %x	十六进制（小写）             fmt.Printf("%x", 42)          → 2a
  %X	十六进制（大写）             fmt.Printf("%X", 42)          → 2A

  %f	浮点数                      fmt.Printf("%f", 3.14)        → 3.140000
  %.2f	浮点数保留 2 位              fmt.Printf("%.2f", 3.1415)   → 3.14
  %e/%E	科学计数法                  fmt.Printf("%e", 1234.5)     → 1.234500e+03
  %.20f 20 位精度（调 bug 用）       fmt.Printf("%.20f", 0.1+0.2) → 0.30000000000000004441

  %t	布尔值                      fmt.Printf("%t", true)        → true

  %s	字符串（不解码转义）          fmt.Printf("%s", "hello")     → hello
  %q	带引号的字符串               fmt.Printf("%q", "hello")     → "hello"
  %q	不可见字符会转义             fmt.Printf("%q", "a\nb")      → "a\nb"

  %c	Unicode 码点对应字符          fmt.Printf("%c", 65)          → A
  %c	中文字符                    fmt.Printf("%c", 22909)       → 好
  %U	Unicode 格式                fmt.Printf("%U", '好')        → U+597D

  %p	指针地址                    fmt.Printf("%p", &v)          → 0xc000012345

  %w	fmt.Errorf 专用，包装错误    fmt.Errorf("xxx: %w", err)
*/

func MathReview() {
	fmt.Println("\n=== 复习：数学运算 ===")

	// 浮点数比较不要用 ==
	// 注意：Go 常量表达式以任意精度求值，0.1+0.2 在编译期可能等于 0.3
	// 要看到真正的精度问题，需要用 float64 变量运算（绕过常量折叠）
	a, b, c := 0.1, 0.2, 0.3
	fConst := a + b // float64 变量相加，真正的 IEEE 754 运算
	fmt.Printf("float64: 0.1+0.2 = %.20f\n", fConst)
	fmt.Printf("float64: 0.3     = %.20f\n", c)
	fmt.Printf("0.1+0.2 == 0.3 ? %t (IEEE 754 精度错误!)\n", fConst == c)
	fmt.Printf("差值: %.20e\n", fConst-c)
	fmt.Printf("math.Abs 容差比较(1e-9): %t\n", math.Abs(fConst-c) < 1e-9)

	// NaN：任何和 NaN 的比较都是 false
	nan := math.NaN()
	fmt.Printf("NaN == NaN ? %t\n", nan == nan) // false! 这是坑

	// 数学函数
	fmt.Printf("math.Pow(2,10)=%f\n", math.Pow(2, 10))
	fmt.Printf("math.Sqrt(16)=%f\n", math.Sqrt(16))
	fmt.Printf("math.Ceil(3.14)=%f\n", math.Ceil(3.14))
	fmt.Printf("math.Floor(3.14)=%f\n", math.Floor(3.14))
}

// ============================================================
// 复习 4：时间和日期
// ============================================================

func TimeReview() {
	fmt.Println("\n=== 复习：时间操作 ===")

	now := time.Now()
	fmt.Printf("当前时间: %v\n", now)
	fmt.Printf("Unix 时间戳: %d\n", now.Unix())

	// 格式化：用参考时间 2006-01-02 15:04:05（Go 特有的时间格式方式）
	fmt.Printf("格式化: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("RFC3339: %s\n", now.Format(time.RFC3339))

	// 时间运算
	future := now.Add(24 * time.Hour)
	past := now.Add(-48 * time.Hour)
	duration := future.Sub(now)
	fmt.Printf("24h 后: %v\n", future)
	fmt.Printf("时间差: %v (%.1f 小时)\n", duration, duration.Hours())

	// 时间解析
	parsed, err := time.Parse("2006-01-02", "2025-01-15") // 同上，使用参考时间
	if err == nil {
		fmt.Printf("解析成功: %v, 星期几: %s\n", parsed, parsed.Weekday())
	}
	_ = past
}
