package review

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================
// 第二阶段 / 模块一：接口 (Interface)
// ============================================================

// 接口即契约：定义"能做什么"，不关心"怎么做"
// Go 的接口是隐式实现的——不需要声明 implements

// --- 示例 1：自定义接口 ---

// Animal 定义一个"能叫"的契约
type Animal interface {
	Speak() string
}

type Dog struct{ Name string }
type Cat struct{ Name string }

func (d Dog) Speak() string { return d.Name + ": 汪汪" }

// Speak 这里是值接收者（只读）
func (c Cat) Speak() string { return c.Name + ": 喵喵" }

// Dog 和 Cat 没有声明 implements Animal，但因为有 Speak() 方法，自动满足接口

// --- 示例 2：空接口 ---

// interface{} (Go 1.18+ 可写 any) 是万能容器，所有类型都满足它

// --- 示例 3：接口组合 ---

// ReadWriter 由两个接口组合而成
type ReadWriter interface {
	Reader // 嵌入 Reader 接口
	Writer // 嵌入 Writer 接口
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// --- 示例 4：标准库接口 ---

// sort.Interface 需要 Len/Less/Swap 三个方法
type Person struct {
	Name string
	Age  int
}
type ByAge []Person

func (p ByAge) Len() int           { return len(p) }
func (p ByAge) Less(i, j int) bool { return p[i].Age < p[j].Age }
func (p ByAge) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// --- 示例 5：指针接收者 vs 值接收者的接口实现 ---

type Greeter interface {
	Greet() string
}

type User struct{ Name string }

func (u User) Greet() string { // 值接收者 → User 和 *User 都实现 Greeter
	return "Hello, " + u.Name
}

type MutableCounter struct{ val int }

func (m *MutableCounter) Greet() string { // 指针接收者 → 只有 *MutableCounter 实现 Greeter
	return fmt.Sprintf("count=%d", m.val)
}

// --- 示例 6：错误接口（你已经在用） ---

// error 就是一个接口:
// type error interface {
//     Error() string
// }

// --- 全部演示 ---

func Interfaces() {
	fmt.Println("=== 新内容：接口 ===")

	// 1. 多态：同一个接口，不同实现
	animals := []Animal{Dog{"旺财"}, Cat{"咪咪"}}
	for _, a := range animals {
		fmt.Println(a.Speak())
	}

	// 2. 空接口 any：可以存任何值
	var anything any
	anything = 42
	fmt.Printf("any 存 int: %v\n", anything)
	anything = "hello"
	fmt.Printf("any 存 string: %v\n", anything)
	anything = Dog{"大黄"}
	fmt.Printf("any 存 struct: %v\n", anything)

	// 3. 类型断言 (comma-ok 安全模式)
	var a Animal = Dog{"旺财"}
	if dog, ok := a.(Dog); ok {
		fmt.Printf("类型断言成功: %s\n", dog.Name)
	}
	if _, ok := a.(Cat); !ok {
		fmt.Println("类型断言失败: 不是 Cat")
	}

	// 4. 类型开关 (type switch)
	checkType := func(v any) {
		switch v := v.(type) { // 注意这个特殊语法 v.(type) 只能用于 switch
		case int:
			fmt.Printf("int: %d\n", v)
		case string:
			fmt.Printf("string: %q, 大写=%q\n", v, strings.ToUpper(v))
		case Animal:
			fmt.Printf("Animal: %s\n", v.Speak())
		default:
			fmt.Printf("未知类型: %T\n", v)
		}
	}
	checkType(42)
	checkType("hello")
	checkType(Dog{"小黑"})
	checkType(3.14)

	// 5. sort.Interface：自定义排序
	people := ByAge{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	fmt.Printf("排序前: %v\n", people)
	sort.Sort(people) // sort.Sort 要求 sort.Interface
	fmt.Printf("按年龄升序: %v\n", people)

	// 6. 值接收者 vs 指针接收者实现接口的差异
	var g Greeter

	g = User{"John"} // ✅ User 有 Greet()
	fmt.Println(g.Greet())
	g = &User{"Jane"} // ✅ *User 也有（自动解引用）
	fmt.Println(g.Greet())

	// g = MutableCounter{10}  // ❌ MutableCounter 没有 Greet()（方法在指针上）
	g = &MutableCounter{10} // ✅ 只有 *MutableCounter 实现了
	fmt.Println(g.Greet())

	// 7. 接口值的零值是 nil（类型和值都是 nil）
	var rdr Reader
	fmt.Printf("nil Reader == nil ? %t\n", rdr == nil) // true

	// 但指针赋给接口后，接口值 != nil（类型信息非 nil）
	fmt.Println("⚠️ 接口陷阱：nil 指针 != nil 接口")
	fmt.Println("   see: https://go.dev/doc/faq#nil_error")
}
