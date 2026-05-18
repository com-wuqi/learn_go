package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func init() {
	fmt.Println("init learn1 ----") // 当被导入时执行
}

type i32 int32

func Learngo1() {
	fmt.Println("Learn1 --")
	const Pi float32 = 3.14159
	fmt.Println(Pi)
	var v1, v2 int
	v1 = 32
	v2 = 31
	fmt.Println(v1, v2)
	v3 := false
	fmt.Println(v3)
	var v4 string
	var v5 [3]int
	fmt.Println(v4, "--", v5)

	fmt.Println(runtime.GOOS, runtime.GOARCH)
	fmt.Println(os.Getenv("PATH"))

	//var logger = log.Default()
	//logger.Println("learngo1 --")
	var v6 = true
	fmt.Println(v6, 1)

	var v7 int8 = 22
	var v8 int16 = 24
	v8 = int16(v7 + int8(v8))
	fmt.Println(v8)

	v9 := rand.Text()
	fmt.Println(v9)
	var v10 i32 = 124
	fmt.Println(v10)

	var ch int = '\u0041'
	var ch2 int = '\u03B2'
	var ch3 int = '\U00101234'
	fmt.Printf("%d - %d - %d\n", ch, ch2, ch3) // integer
	fmt.Printf("%c - %c - %c\n", ch, ch2, ch3) // character
	fmt.Printf("%X - %X - %X\n", ch, ch2, ch3) // UTF-8 bytes
	fmt.Printf("%U - %U - %U", ch, ch2, ch3)   // UTF-8 code point

	//if unicode.IsUpper('A') {
	//	fmt.Println("12")
	//}

	var v11 string
	v11 += "a string"

}

func Learngo2() {
	fmt.Println("Learngo2 --")
	var v1 string
	var v2 string
	v1 = "a b d e f"
	v2 = "0123456789 34 34"
	if strings.HasPrefix(v1, "a") {
		fmt.Println("前缀")
	}
	if strings.HasSuffix(v1, "e f") {
		fmt.Println("后缀")
	}
	if strings.Contains(v1, "b d e") {
		fmt.Println("包含")
	}
	// Index 返回字符串 str 在字符串 s 中的索引（str 的第一个字符的索引），-1 表示字符串 s 不包含字符串 str
	fmt.Println(strings.Index(v2, "34"))     // -> 3, first
	fmt.Println(strings.LastIndex(v2, "34")) // -> 14

	v4 := "0123456abcd"
	v4 = strings.Replace(v4, "abcd", "789", 1)
	// Replace 用于将字符串 str 中的前 n 个字符串 old 替换为字符串 new，并返回一个新的字符串，
	// 如果 n = -1 则替换所有字符串 old 为字符串 new：
	fmt.Println(v4)

	v5 := "abc abc abc abcd"
	v6 := strings.Count(v5, "abc")
	fmt.Println("v6=", v6)

	v7 := strings.Split(v2, " ") // 切片
	fmt.Println(v7, "\n", v7[0])

	// strings.Repeat(s, count int) string // Repeat 用于重复 count 次字符串 s 并返回一个新的字符串
	// strings.ToLower(s) string // 全小写
	fmt.Println("--")
	v8 := "\n123 \n  \n "
	v8 = strings.TrimSpace(v8) // 切除空白
	fmt.Println(v8)
	fmt.Println("--")
	v9 := "cut123cut"
	v9 = strings.TrimRight(v9, "cut")
	fmt.Println(v9)
	v9 = strings.TrimLeft(v9, "cut")
	fmt.Println(v9)

	v10 := " \n12 13 14\n15 \n16\n 17"
	var v11 []string
	v11 = strings.Fields(v10) // 可以，但是考虑首先 strings.TrimSpace
	fmt.Println(v11)
	v12 := "11.12.13.14"
	fmt.Println(strings.Split(v12, "")) // [1 1 . 1 2 . 1 3 . 1 4]
	// 每个字符切开
	fmt.Println(strings.Split(v12, "...")) // [11.12.13.14]
	// 一个没切，做为一个元素

}
