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
