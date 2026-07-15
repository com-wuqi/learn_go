package practice

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ============================================================
// 反射练习（学完 review/reflect.go 后做）
// ============================================================

// ============================================================
// 练习 31：StructToMap
// 用反射将结构体转为 map[string]interface{}，key 为字段名
// ============================================================

// [TODO] 将任意结构体转为 map，只处理导出的顶层字段（非嵌套）
// 示例: StructToMap(User{Name:"Alice", Age:30}) → map["Name":"Alice", "Age":30]
// 提示: reflect.TypeOf().NumField(), val.Field(i).Interface()
func StructToMap(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	val := reflect.ValueOf(s)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		m[typ.Field(i).Name] = val.Field(i).Interface()
	}
	return m
}

// ============================================================
// 练习 32：FillDefaults
// 用反射读取结构体的 `default` tag，将零值字段填充为默认值
// ============================================================

// [TODO] 传入结构体指针，对零值字段填充 default tag 指定的值
// 示例: &User{Name:"", Age:0} 有 tag `default:"anonymous"` → Name 变为 "anonymous"
// 提示: 只处理 string 类型的零值字段（用 CanSet() 检查、SetString() 赋值）
func FillDefaults(ptr interface{}) {
	val := reflect.ValueOf(ptr).Elem() // 类似 *
	typ := reflect.TypeOf(ptr).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanSet() && field.Kind() == reflect.String && field.String() == "" {
			defaultVal := typ.Field(i).Tag.Get("default")
			field.SetString(defaultVal)
		}
	}
}

// ============================================================
// 练习 33：PrintTags
// 用反射遍历结构体字段，打印 json 和 default tag
// ============================================================

// [TODO] 输出格式: "字段名: json=xxx, default=xxx"
// 提示: field.Tag.Get("json")
func PrintTags(s interface{}) string {
	builder := strings.Builder{}
	val := reflect.ValueOf(s) // 不是指针，不需要 Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		//field := val.Field(i)
		typeField := typ.Field(i)
		defaultTag := typeField.Tag.Get("default")
		jsonTag := typeField.Tag.Get("json")
		builder.WriteString(fmt.Sprintf("%s: json=%s, default=%s\n", typeField.Name, jsonTag, defaultTag))
	}
	return builder.String()
}

// 辅助：练习 31-33 使用的示例结构体
type ReflectUser struct {
	Name string `json:"name" default:"anonymous"`
	Age  int    `json:"age"`
	City string `json:"city" default:"unknown"`
}

// ============================================================
// 练习 36：panic/recover — 模拟服务降级
// ============================================================

// [TODO] SafeCall 执行 fn, 如果 fn panic 则 recover 并返回 panic 的信息作为 error
// 返回 nil 表示执行成功（无 panic）
// 提示: defer + recover(), fmt.Errorf("panic: %v", r)
func SafeCall(fn func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	fn()
	return nil
}
