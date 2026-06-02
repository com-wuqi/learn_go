package utils

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func init() {
	fmt.Println()
}

func LearnGo2() {

}

func stringToInt(origin string) (int, error) {
	anInt, anError := strconv.Atoi(origin)
	return anInt, anError
}

func mySqrt(f float64) (v float64, ok bool) {
	if f < 0 {
		return math.NaN(), false
	}
	return math.Sqrt(f), true
}

func mySqrt2(f float64) (v float64, ok bool) {
	// 感觉比上面的更好一些
	v1 := math.Sqrt(f)
	if v1 == math.NaN() {
		return math.NaN(), false
	} else {
		return v1, true
	}
}

func Learngo4() {
	fmt.Println("Learngo4 --")
	v1 := true
	var v2 = false
	fmt.Println(v1, v2)
	if !(v1 && v2) {
		fmt.Println("v1 and v2 is false")
	}
	if v1 || v2 {
		fmt.Println("v1 and v2 is true")
	}
	// useful examples
	var v3 = "hello"
	var v4 = ""
	var v5 = &v4
	if len(v3) > 0 && len(v4) > 0 {
		fmt.Println("v3 and v4 is ", v3, v4)
	}
	v4 = "world"
	if len(v3) > 0 && len(v4) > 0 {
		fmt.Println("v3 and v4 is ", v3, *v5)
	}
	var v6 int8
	if v6 = 10; math.Abs(float64(v6)) > 0.1 {
		fmt.Println("v6 is ", v6)
	}
	v7 := "1234567890"
	v7Int, v7Error := stringToInt(v7)
	if v7Error != nil {
		fmt.Println("v7 is not an Int ", v7Error)
		os.Exit(1)
	} else {
		fmt.Println("v7 is ", v7Int)
	}
	fmt.Println("do sth else") // 标记 os.Exit()

	// 简单写法
	//value, err := pack1.Function1(param1)
	//if err != nil {
	//	fmt.Printf("An error occured in pack1.Function1 with parameter %v", param1)
	//	return err
	//}

	// 简洁一点, 但是 value 后续不可用
	if value, err := stringToInt(v7); err == nil {
		fmt.Println(value) // 1234567890
	}
	////fmt.Println(value) // wrong

	printlnCount, printlnError := fmt.Println(1234) // 1234
	// 当打印到控制台时，可以将该函数返回的错误忽略；
	// 但当输出到文件流、网络流等具有不确定因素的输出对象时，
	// 应该始终检查是否有错误发生
	if printlnError != nil {
		fmt.Println("printlnError:", printlnError)
		os.Exit(1)
	} else {
		fmt.Println(printlnCount) // 计数
	}
}

func LearnGo5() {
	fmt.Println("LearnGo5 --")
	var v1 = 12
	var v2 = 14
	switch v2 {
	case 14:
		fmt.Println("v2 is 14")
		fallthrough
	case 15:
		fmt.Println("v2 is 15")
		// 不成立, 仍然结束
	default:
		fmt.Println("v2 is default")
	}
	switch {
	case v1 == 12:
		fmt.Println("v1 is 12")
		fallthrough // 继续判断
	case v1 >= 12: // 仍然成立
		fmt.Println("v1 is through 12")
		// 不需要 break
	default:
		fmt.Println("default")
	}

	// for
	for v3 := 0; v3 < 5; v3++ {
		fmt.Println("v3 is", v3)
	}
	// 5.4
}
