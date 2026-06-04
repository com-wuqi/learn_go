package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
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

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday: // 多个参数
		{
			fmt.Println("it is a weekend")
		}
	default:
		fmt.Println("it is a weekday")
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

	whatAmI := func(i any) {
		switch i.(type) {
		case bool:
			fmt.Println(i, "I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Println("default")
		}
	}
	whatAmI(true)
	/*
		处理 JSON 反序列化：解析未知结构的 JSON 时，数据通常被解析为 map[string]any 或 []any，
		必须用 Type Switch 来安全地提取值。
		错误处理：判断一个 error 接口底层具体是哪种自定义错误类型
		例如 switch err := err.(type) { case *MyCustomError: ... }
	*/
	// 类似 fmt 包的功能：fmt.Println
	// 需要接收任何类型，并根据它是字符串、数字还是切片来执行完全不同的格式化逻辑
	// 例如:
	HandleLog("a string")
	HandleLog(true)
	HandleLog(fmt.Errorf("database connection failed")) // 构建一个错误
	HandleLog(float64(1234.1))                          // default

	// 需要注意的是, 对于不同类型但是处理逻辑相同, 应该使用泛型, 例如
	fmt.Println("sum up:", sumSomeNumbers(1, 2, 3, 4, 5))
	fmt.Println("sum up:", sumSomeNumbers(1.1, 1.2, 1.4, 1.5))
	// 但这并不意味着你可以在同一次函数调用中混合传入 多个不同的类型
	// fmt.Println("sum up:", sumSomeNumbers("a","b", 3)) // 无法推断T

	/*
		在 Go 中，像 1、1.1、2 这样直接写出来的数字字面量，被称为无类型常量。
		它们在没有被赋值给具体变量之前，没有固定的类型（既不是严格的 int，也不是严格的 float64），
		而是具有“灵活性”
	*/

	// for
	for v3 := 0; v3 < 5; v3++ {
		fmt.Println("v3 is", v3)
	}

	for {
		fmt.Println("loop and break")
		break
	}

	for v4 := 3; v4 >= 0; {
		v4 = v4 - 1
		if v4 == 1 {
			continue
		}
		fmt.Printf("The variable v4 is now: %d\n", v4)
	}
	// 标签的名称是大小写敏感的，为了提升可读性，一般建议使用全部大写字母

LABEL1:
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 4 || i == 4 {
				continue LABEL1
			}
			fmt.Printf("i is: %d, and j is: %d\n", i, j)
		}
	}
}

func HandleLog(msg any) {
	switch i := msg.(type) {
	case string, int, float32, bool:
		fmt.Println("[info]", i)
	case error:
		fmt.Println("[error]", i.Error())
	case []byte:
		{
			var aJson map[string]any
			if err := json.Unmarshal(i, &aJson); err == nil {
				fmt.Println("[json]", aJson)
			} else {
				fmt.Println("[byte]", string(i))
			}
		}
	default:
		fmt.Printf("[default] Type: %T, Value: %+v\n", i, i)
	}
}

func sumSomeNumbers[T int | float32 | float64 | string](items ...T) T {
	var sum T
	for _, number := range items {
		sum += number
	}
	return sum
}
