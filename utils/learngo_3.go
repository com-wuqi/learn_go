package utils

import (
	"fmt"
	"slices"
)

func LearnGo6() {
	fmt.Println("LearnGo6")
	// 函数
	// 函数重载（function overloading）指的是可以编写多个同名函数，
	// 只要它们拥有不同的形参 / 或者不同的返回值，
	// 在 Go 里面函数重载是不被允许的。
	// 这将导致一个编译错误：funcName redeclared in this book, previous declaration at lineno
	// Go 语言目前不支持函数重载（Function Overloading），并且 Go 语言的设计团队明确表示，未来也没有计划加入这一特性。
	// 6.1

	// https://gobyexample-cn.github.io/functions
	// https://gobyexample-cn.github.io/multiple-return-values
	// https://gobyexample-cn.github.io/variadic-functions
	// 在函数调用时，像切片（slice）、字典（map）、接口（interface）、通道（channel）这样的引用类型都是默认使用引用传递（即使没有显式的指出指针）。
	// 指针也是变量类型，有自己的地址和值，通常指针的值指向一个变量的地址。所以，按引用传递也是按值传递。

	fmt.Println("v61")
	_, v1 := v61(12, 13) // 这里的 _ 和py3一样，用于丢弃
	fmt.Println("v61:", v1)

	/* 在函数中改变外部变量
	 */

	// 1. 只是返回值，由调用者处理
	// 2. 闭包
	v2 := 1
	modefyV2 := func(a int) {
		v2 += a
	}
	modefyV2(12) // 外部使用
	fmt.Println("v2:", v2)
	// 3. 利用“引用语义”的内置类型（Slice, Map, Channel）
	// 当你把它们作为参数传递时，虽然拷贝了它们的“描述符（Header）”，
	// 但底层指向的是同一块内存。因此，修改它们内部的元素，会影响外部。
	v3 := map[string]int{"a": 1, "b": 2}
	changeMap(v3) // v3: map[a:1 added:1234 b:2]
	fmt.Println("v3:", v3)
	// 对于切片([]int{1,2,3})，如果你在函数内部使用 append 导致底层数组扩容，
	// 或者直接用 s = []int{...} 重新赋值，这不会影响外部！
	// 因为此时你修改的是“描述符”本身（值传递的拷贝），而不是底层数据。
	// 4. 结构体指针与“指针接收者”方法
	v4 := Example1{inSide: 1}
	changeExample1(&v4)
	fmt.Println("v4:", v4)

}

type Example1 struct {
	inSide int
}

func changeExample1(v1 *Example1) {
	v1.inSide = 2
}

func changeMap(v1 map[string]int) {
	v1["added"] = 1234
}

func v61(a int, b int) (c int, d int) {
	// 尽量使用命名返回值：会使代码更清晰、更简短，同时更加容易读懂。
	if a < b {
		c, d = a, b
	} else {
		c, d = b, a
	}
	//return d,c // 仍然可以无视
	return
}

func LearnGo7() {
	fmt.Println("LearnGo with example")
	var v1 [5]int
	v1[0] = 1
	v2 := [...]string{"1", "2", "3"} // 这里长度已经确定
	fmt.Println(v2)
	//v2[10] = 11
	var v3 [3]string
	v3[0] = "hello"
	v3[1] = "world"
	v3[2] = "!"
	v4 := make([]string, 1) // 内部有一个 "" 了
	// 长度是动态的,分配到堆上, GC回收, 更加灵活
	v4 = append(v4, "hello")
	v4 = append(v4, "world")
	v4 = append(v4, "!")
	fmt.Println(v4)
	// 这里经历了三次扩容, [ + "" +   + hello +   + world +   + ! + ]
	// var v4 []string 可以解决
	// 或者 v4 := make([]string, 0)
	v5 := make([]int, 0, 2) // 提前给定长度可以避免反复扩容
	v5 = append(v5, 1)
	v5 = append(v5, 2)
	v5 = append(v5, 3)
	v5 = append(v5, 4)
	fmt.Println(v5)
	// https://go.dev/blog/slices-intro
	fmt.Println(v5[0:1]) // 支持常见切片操作
	fmt.Println(v5[1:])
	//fmt.Println(v5[-1:]) // invalid argument: index -1 (constant of type int) must not be negative
	fmt.Println(v5[:1:1])
	//copy(v4, v2) // invalid copy: argument must be a slice; have v2 (variable of type [3]string)
	// 切片到切片可以直接copy,数组和数组可以 = , 数组到切片需要类型转换
	copy(v4, v2[:])
	fmt.Println(v4) // [1 2 3 !]
	v51 := []string{"hello", "world"}
	v52 := []string{"a"}
	copy(v52, v51) // 这里长度不足, 发生截断
	fmt.Println(v52)
	v6 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	copy(v6[:4], v6[5:]) // 6~10 复制到 1~5 并静默截断
	fmt.Println(v6)
	// 安全处理“内存重叠”
	slices.Reverse(v6) // 颠倒 v6
	fmt.Println(v6)

	// 因为 Go 的字符串是 UTF-8 编码的，
	// 如果包含中文或 Emoji，一个字符可能占用多个字节。
	// 直接反转 []byte 会把中文字符截断，变成乱码。
	str1 := "hello world!"
	runes := []rune(str1) // 先转成 []rune (Unicode 字符切片)
	slices.Reverse(runes)
	fmt.Println(string(runes))
	// https://gobyexample-cn.github.io/slices
}

func LearnGo8() {
	fmt.Println("LearnGo8 with example")
	var v1 map[string]any
	v1 = make(map[string]any)
	v1["hello"] = "world"
	fmt.Println(v1)
	v1["a"] = 1
	v1["b"] = 2
	v1["c"] = 3
	fmt.Println(v1)      // 打印整个map
	fmt.Println(len(v1)) // map大小
	delete(v1, "a")
	fmt.Println(v1)
	v1["b"] = 1 // 修改元素
	fmt.Println(v1)
	// 更加安全的访问
	_, _ok := v1["unknown"]
	fmt.Println("v1[\"unknown\"] ok:", _ok)

	// Range
	// https://gobyexample-cn.github.io/range
	v2 := [...]int{1, 2, 3, 4}
	for _, v := range v2 {
		fmt.Print(v, " ")
	}
	fmt.Println()
	for k, v := range v2 { // index-value
		fmt.Print(k, ": ", v, " | ")
	}
	fmt.Println()
	v3 := map[string]any{"a": 1, "b": 2}
	fmt.Println(v3)
	for k, v := range v3 { // key-v, 和上文类似
		fmt.Print(k, ": ", v, " | ")
	}
	fmt.Println()

	// 另附, 一点额外的想法
	// v5 := [...]map[string]any{make(map[string]string), make(map[string]int), make(map[string]float32)}
	// 后面的只能是any
	// 在 Go 语言中，map[string]string 和 map[string]any 是两种完全不同的类型，Go 不允许将它们直接互相赋值或隐式转换。
	v5 := []any{
		map[string]int{"a": 1, "b": 2},
		//map[int]int, // 这样不行， 需要make
		make(map[int]int),
		make(map[string]int),
	}
	fmt.Println(v5)
	v6 := map[string]int{"a": 1, "b": 2} // 希望转成 string-any
	v61 := make(map[string]any)
	// 我们只能手动遍历
	for k, v := range v6 {
		v61[k] = v // 自动装箱
	}
	fmt.Println("v61", v61)
}

func LearnGo9() {
	fmt.Println("LearnGo9")
	v1 := []string{"hello", "world"}
	fmt.Println(v1)
	learnGoExample(v1...)

	// defer return前执行，panic兜底，多路径返回支持，甚至修改返回值
	// 立刻求值（defer当前行），后进先出
	_, err := learnDefer("", "", "")
	if err != nil {
		fmt.Println(err)
	}
	// 后进先出允许对资源按序加锁
	// 一点闭包
	v2 := aCounter()
	v2()
	v2()
	fmt.Println(v2()) // 3

	// https://learnku.com/docs/the-way-to-go/application-closure-function-as-a-return-value/3607
	// 工厂函数/函数式选项模式 的构造函数 https://chat.qwen.ai/c/79ff671a-42fe-4f6b-8b83-b0241e04d237

}

func learnGoExample(v1 ...string) (v2 []string) {
	// 这里看起来同样视为声明
	fmt.Println("len()", len(v1))
	// 必须参数数目大于0时考虑判断，for会被静默跳过
	for _, _v1 := range v1 {
		v2 = append(v2, _v1)
	}
	//var v3 string
	//strings.Join(v1,v3) // v3 is a slice
	//fmt.Println(v3)
	return v2
}

func learnDefer(v1 ...string) (string, error) {
	defer fmt.Println("deferred 1")
	defer fmt.Println("deferred 2")
	defer fmt.Println("deferred 3")
	fmt.Println("run 1")
	// IIFE
	func() {
		defer fmt.Println("deferred 4")
		fmt.Println("run 2")
	}()
	// 1 2 4 3 2 1
	// 赋值
	v2 := func() {
		defer fmt.Println("deferred 5")
		fmt.Println("run 3")
	}
	v2()
	/*
		run 1
		run 2
		deferred 4
		run 3
		deferred 5
		deferred 3
		deferred 2
		deferred 1
	*/
	return v1[0], nil
}

func aCounter() func() int {
	v1 := 0
	// 这个函数返回的函数使用的v1是在heap上分配的，这个地址指向的变量仍然存在
	// 换言之，这个函数和它的环境一起被保留, 环境中的变量可以被修改
	return func() int {
		v1++
		return v1
	}
}
