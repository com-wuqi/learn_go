package review

import (
	"encoding/json"
	"fmt"
	"os"
)

// ============================================================
// 补充模块：文件 I/O 与 JSON 编解码
// ============================================================
// 目标：补上 FileStore 和 JSON tag 练习所需的基础

// --- 示例 1：os.ReadFile / os.WriteFile（Go 1.16+） ---

func FileDemo() {
	path := "/tmp/learngo_demo.txt"

	// 写入文件：os.WriteFile(路径, 内容, 权限)
	os.WriteFile(path, []byte("hello, 文件IO\n"), 0644)
	// 0644 = rw-r--r--（owner 可读写，其他人只读）

	// 读取文件：os.ReadFile 返回 []byte 和 error
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读文件失败:", err)
		return
	}
	fmt.Printf("读取内容: %s", data)

	// 清理
	os.Remove(path)
}

// --- 示例 2：JSON 序列化（结构体 → JSON 字节） ---

type PersonJSON struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func JSONDemo() {
	// Marshal：结构体 → JSON
	p := PersonJSON{Name: "Alice", Age: 30}
	b, _ := json.Marshal(p)
	fmt.Printf("json.Marshal: %s\n", b)
	// 输出: {"name":"Alice","age":30}
	// 注意：字段名变成 snake_case 是因为 tag 写了 `json:"name"`

	// MarshalIndent：带缩进，方便阅读
	b2, _ := json.MarshalIndent(p, "", "  ")
	fmt.Printf("json.MarshalIndent:\n%s\n", b2)
	// 输出:
	// {
	//   "name": "Alice",
	//   "age": 30
	// }

	// Unmarshal：JSON → 结构体
	jsonStr := `{"name":"Bob","age":25}`
	var p2 PersonJSON
	json.Unmarshal([]byte(jsonStr), &p2) // 注意传指针！
	fmt.Printf("json.Unmarshal: %+v\n", p2)
}

// --- 示例 3：文件 + JSON 组合 = FileStore 的核心逻辑 ---

func FileJSONDemo() {
	path := "/tmp/learngo_store.json"
	type Store struct {
		Data map[string]string `json:"data"`
	}

	// 写入：构造结构体 → Marshal → WriteFile
	s1 := Store{Data: map[string]string{"key1": "val1"}}
	b, _ := json.MarshalIndent(s1, "", "  ")
	os.WriteFile(path, b, 0644)

	// 读取：ReadFile → Unmarshal → 结构体
	raw, _ := os.ReadFile(path)
	var s2 Store
	json.Unmarshal(raw, &s2)
	fmt.Printf("恢复的数据: %v\n", s2.Data)

	os.Remove(path)
}

// --- 示例 4：defer 关闭文件 ---
// 大文件场景用 os.Open + defer f.Close()，ReadFile/WriteFile 不需要 defer

func DeferFileDemo() {
	f, err := os.Create("/tmp/learngo_defer.txt")
	if err != nil {
		return
	}
	defer f.Close() // 函数返回时自动关闭

	f.WriteString("defer 保证文件被关闭\n")
	// 即使后面 panic/recover，defer 也会执行
}
