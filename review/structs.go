package review

import (
	"encoding/json"
	"fmt"
)

// ============================================================
// 第二阶段 / 模块二：结构体与方法 (Ch10 核心)
// ============================================================
// 你已掌握：值/指针接收者、String() 接口、工厂函数
// 新内容：结构体嵌入（组合 > 继承）、结构体标签(tag)

// --- 示例 1：结构体嵌入（匿名字段） ---
// Go 没有继承，用"嵌入"实现代码复用 — 对分布式来说，配置组合无处不在

type BaseConfig struct {
	Host string
	Port int
}

// ServiceConfig 嵌入 BaseConfig（匿名字段），直接获得 Host/Port
// 这叫"提升"（promotion）
type ServiceConfig struct {
	BaseConfig        // 匿名字段 — ServiceConfig 可直接访问 Host、Port
	Name       string `json:"name"` // ← 结构体标签（见示例 3）
}

// --- 示例 2：嵌入 vs 具名字段 ---

type EmbeddedExample struct {
	BaseConfig // 嵌入：可直接 s.Host
}

type HasAExample struct {
	base BaseConfig // 具名字段：必须 s.base.Host（无法省略）
}

// --- 示例 3：结构体标签 (Struct Tags) ---
// tag 是附加在字段上的字符串，JSON/protobuf/ORM 等库通过反射读取 tag

type ServerConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	MaxConns int    `json:"max_conns,omitempty" yaml:"max_conns"`
	internal string // 小写字母开头 = 不导出，json.Marshal 会忽略
}

// --- 示例 4：方法冲突 ---
// 如果嵌入的两个类型有同名方法，编译器会报错，必须显式解决

type A struct{}
type B struct{}

func (a A) String() string { return "A" }
func (b B) String() string { return "B" }

type AB struct {
	A
	B
}

// AB 同时继承了 A.String() 和 B.String() — 不加处理直接调用会编译报错

// --- 全部演示 ---

func Structs() {
	fmt.Println("=== 结构体与方法（嵌入 & tag）===")

	// 1. 嵌入：直接访问提升的字段
	svc := ServiceConfig{
		BaseConfig: BaseConfig{Host: "localhost", Port: 8080},
		Name:       "my-service",
	}
	fmt.Printf("嵌入: %s:%d (%s)\n", svc.Host, svc.Port, svc.Name)
	// svc.Host 等同于 svc.BaseConfig.Host（提升）

	// 2. 具名字段 vs 嵌入
	hasA := HasAExample{base: BaseConfig{Host: "db.local", Port: 5432}}
	fmt.Printf("具名字段: %s:%d\n", hasA.base.Host, hasA.base.Port)
	// hasA.Host 编译报错 — 不提升

	// 3. 结构体标签：json.Marshal 会读取
	cfg := ServerConfig{Host: "0.0.0.0", Port: 9090, MaxConns: 100, internal: "secret"}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("JSON 序列化（tag 决定字段名）:\n%s\n", b)
	// 注意: internal 小写开头，不导出，不会出现在 JSON 中

	// MaxConns 标签包含 omitempty — 零值时不输出
	cfg2 := ServerConfig{Host: "0.0.0.0", Port: 9090} // MaxConns=0
	b2, _ := json.MarshalIndent(cfg2, "", "  ")
	fmt.Printf("omitempty 效果 (MaxConns=0):\n%s\n", b2)

	// 4. 嵌入 + 方法：子类型自动获得父类型的方法
	// （如果有同名方法冲突，参考上面的 AB 结构体）
	fmt.Println("\n嵌入方法：如果 BaseConfig 有 String()，ServiceConfig 也会继承")
}
