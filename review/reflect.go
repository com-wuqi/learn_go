package review

import (
	"fmt"
	"reflect"
)

// ============================================================
// 模块：反射 (reflect) — Ch11.10
// ============================================================
// 反射让程序在运行时检查类型和值。ORM、RPC、JSON 编解码都靠它。
// 代价：代码复杂、性能差、编译期类型安全丢失。谨慎使用。

// --- 示例 1：reflect.TypeOf 和 reflect.ValueOf ---

func ReflectBasics() {
	fmt.Println("=== 反射基础 ===")

	x := 42
	t := reflect.TypeOf(x)  // 获取类型信息
	v := reflect.ValueOf(x) // 获取值信息

	fmt.Printf("TypeOf: %v (%s)\n", t, t.Kind()) // int
	fmt.Printf("ValueOf: %v\n", v)               // 42
	fmt.Printf("CanInt: %t, Int=%d\n", v.CanInt(), v.Int())
}

// --- 示例 2：检查结构体字段 ---

type DemoUser struct {
	Name string `json:"name" default:"anonymous"`
	Age  int    `json:"age"`
}

func ReflectStruct() {
	u := DemoUser{Name: "Alice", Age: 30}

	// 值对象：可以读，不能修改
	val := reflect.ValueOf(u)
	typ := val.Type()

	fmt.Printf("\n=== 结构体反射: %s ===\n", typ.Name())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)        // reflect.StructField（类型信息）
		value := val.Field(i)        // reflect.Value（值）
		tag := field.Tag.Get("json") // 读取 json tag

		fmt.Printf("  %s (%s) = %v  [json:\"%s\"]\n",
			field.Name, field.Type, value, tag)
	}
}

// --- 示例 3：修改结构体字段（需要传指针） ---

func ReflectModify() {
	u := DemoUser{Name: "Alice", Age: 30}
	v := reflect.ValueOf(&u).Elem() // &u 指向 User，.Elem() 取出指针指向的值

	// 修改 Name 字段（必须检查 CanSet）
	nameField := v.FieldByName("Name")
	if nameField.CanSet() {
		nameField.SetString("Bob") // 修改成功
	}

	fmt.Printf("\n=== 反射修改结构体 ===\n")
	fmt.Printf("修改后: %+v\n", u) // {Name:Bob Age:30}

	// ⚠️ 如果传的是 u 而不是 &u，CanSet() 返回 false，SetString 会 panic
}

// --- 示例 4：读取 tag 中的自定义信息 ---

func ReflectTags() {
	u := DemoUser{}
	typ := reflect.TypeOf(u)

	fmt.Printf("\n=== 反射读取 Tag ===\n")
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		defaultVal := field.Tag.Get("default")
		fmt.Printf("  %s: json=%q, default=%q\n", field.Name, jsonTag, defaultVal)
	}
}

// --- 示例 5：Kind 判断（类型家族） ---

func ReflectKind() {
	fmt.Printf("\n=== Kind 类型家族 ===\n")

	check := func(v any) {
		kind := reflect.TypeOf(v).Kind()
		fmt.Printf("  %T → Kind=%s", v, kind)
		if kind == reflect.Struct {
			fmt.Print(" （是结构体）")
		}
		fmt.Println()
	}

	check(42)               // Kind=Int
	check("hello")          // Kind=String
	check(DemoUser{})       // Kind=Struct
	check(&DemoUser{})      // Kind=Ptr
	check([]int{1, 2, 3})   // Kind=Slice
	check(map[string]int{}) // Kind=Map
}

// --- 示例 6：reflect 的常见坑 ---

func ReflectPitfalls() {
	fmt.Println("\n=== 反射常见坑 ===")

	// 坑1: 非导出字段（小写开头）不可 Set
	type secret struct{ x int }
	s := secret{x: 1}
	f := reflect.ValueOf(&s).Elem().Field(0)
	fmt.Printf("非导出字段 CanSet: %t（小写字段无法反射修改）\n", f.CanSet())

	// 坑2: 传值 vs 传指针
	val := reflect.ValueOf(User{})
	fmt.Printf("ValueOf(struct): CanSet=%t（值不可改）\n", val.Field(0).CanSet())

	ptr := reflect.ValueOf(&User{}).Elem()
	fmt.Printf("ValueOf(&struct).Elem(): CanSet=%t（指针可改）\n", ptr.Field(0).CanSet())
}
