package practice

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// ============================================================
// 接口练习（学完 review/interfaces.go 后做）
// ============================================================

// 练习 15：实现 sort.Interface 对字符串按长度排序
type ByLen []string

func (b ByLen) Len() int           { return len(b) }
func (b ByLen) Less(i, j int) bool { return len(b[i]) < len(b[j]) }
func (b ByLen) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// 练习 16：KVStore 接口，分别用内存 map 和文件实现
type KVStore interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type MemStore struct {
	data map[string]string
}

func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]string)}
}
func (m *MemStore) Get(key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}
func (m *MemStore) Set(key string, value string) error {
	m.data[key] = value
	return nil
}
func (m *MemStore) Delete(key string) error {
	if _, ok := m.data[key]; ok {
		delete(m.data, key)
		return nil
	}
	return errors.New("key not found")
}

type FileStore struct {
	path string
	data map[string]string
}

func NewFileStore(path string) *FileStore {
	b, err := os.ReadFile(path)
	if err != nil {
		return &FileStore{path: path, data: make(map[string]string)}
	}
	var data map[string]string
	err = json.Unmarshal(b, &data)
	if err != nil {
		return &FileStore{path: path, data: make(map[string]string)}
	}
	return &FileStore{path: path, data: data}
}
func (f *FileStore) Get(key string) (string, error) {
	if val, ok := f.data[key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}
func (f *FileStore) Set(key string, value string) error {
	f.data[key] = value
	return f.Save()
}
func (f *FileStore) Delete(key string) error {
	if _, ok := f.data[key]; ok {
		delete(f.data, key)
		return f.Save()
	}
	return errors.New("key not found")
}
func (f *FileStore) Save() error {
	b, err := json.MarshalIndent(f.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, b, 0644)
}

// 练习 17：类型断言 — 遍历 Animal 切片
type AnimalI interface {
	Speak() string
}

type DogI struct{ Name string }
type CatI struct{ Name string }

func (d DogI) Speak() string { return d.Name + ": 汪汪" }
func (c CatI) Speak() string { return c.Name + ": 喵喵" }

func DescribeAnimals(animals []AnimalI) []string {
	var result []string
	for _, animal := range animals {
		if animal != nil {
			result = append(result, animal.Speak())
		}
	}
	return result
}

// 练习 18：Plugin 模式 — Pipeline 串联处理
type Plugin interface {
	Process(s string) string
}

type UppercasePlugin struct{}

func (u UppercasePlugin) Process(s string) string {
	return strings.ToUpper(s)
}

type ReversePlugin struct{}

func (r ReversePlugin) Process(s string) string {
	var runes []rune
	for _, r := range s {
		runes = append(runes, r)
	}
	slices.Reverse(runes)
	return string(runes)
}

func RunPipeline(input string, plugins []Plugin) string {
	result := input
	for _, plugin := range plugins {
		result = plugin.Process(result)
	}
	return result
}

// 练习 19：自定义错误类型（实现 error 接口）
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Msg)
}

// 练习 20：实现 fmt.Stringer 接口
type IPAddr [4]byte

func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// 练习 21：nil 接口陷阱
func IsNil(v interface{}) bool {
	return v == nil
}

func NilInterfaceDemo() {
	var p *int = nil
	var v interface{} = p
	fmt.Printf("p == nil: %t\n", p == nil) // true
	fmt.Printf("v == nil: %t\n", v == nil) // false
	fmt.Printf("IsNil(p): %t\n", IsNil(p)) // true
	fmt.Printf("IsNil(v): %t\n", IsNil(v)) // false（为什么?）
}

// 练习 22：接口组合
type Printer interface {
	Print() string
}

type Scanner interface {
	Scan() string
}

type AllInOne interface {
	Printer
	Scanner
}

type MyDevice struct{}

func (d MyDevice) Print() string { return "" }
func (d MyDevice) Scan() string  { return "" }
