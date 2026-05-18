package main

import (
	"LearnGo/utils"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func swap(s string, i string) (string, string) {
	return i, s
}

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		fmt.Println("i =", 100/i)
	}
	var a = [...]int{2, 4, 3}
	fmt.Println(a)
	var b = []int{1, 2, 3}
	fmt.Println(b)
	var c [3]int // [0 0 0]
	fmt.Println(c)
	fmt.Println(len(c), c[0]) // 长度，切片
	var d = &a
	for i := range c {
		c[i]++ // 依次累加
	}
	fmt.Println(c)
	for i, v := range d { // i,v 均可选
		fmt.Println(i, v)
	}
	for i := 0; i < len(c); i++ {
		fmt.Println(c[i]) // 1 1 1
	}
	var v1 [2]image.Point
	v1[0].X = 1
	v1[0].Y = 1
	v1[1] = image.Point{X: 2, Y: 4}
	fmt.Println(v1)
	fmt.Printf("%d\n", a[0])

	var decoder2 = [...]func(r io.Reader) (image.Image, error){
		png.Decode,
		jpeg.Decode,
	}
	decoder2[0] = func(r io.Reader) (image.Image, error) {
		fmt.Println("sth here")
		return jpeg.Decode(r)
	}
	fmt.Println(decoder2)
	//done := make(chan struct{})
	//for i := 0; i < 10; i++ {
	//	go func(id int) {
	//		<-done // 阻塞直到 done 被关闭
	//		fmt.Println(id, "exiting")
	//	}(i)
	//}
	//time.Sleep(1 * time.Second)
	//close(done)
	v2 := "hello world"
	fmt.Println(v2)
	// func FunctionName (a typea, b typeb) (t1 type1, t2 type2)
	fmt.Println(swap("hello", "world"))
	utils.Learngo1()
	utils.Learngo2()

}
