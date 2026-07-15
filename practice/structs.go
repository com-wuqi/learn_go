package practice

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ============================================================
// 结构体与方法练习（学完 review/structs.go 后做）
// ============================================================

// 练习 24：结构体嵌入
// Database 嵌入 Config（匿名字段），直接提升 Host/Port
type Config struct {
	Host string
	Port int
}

type Database struct {
	// TODO: 嵌入 Config，再添加 Name 字段
	Config
	Name string
}

// TODO: Connect 返回 "postgres://Host:Port/Name"
func (d *Database) Connect() string {
	return fmt.Sprintf("postgres://%s:%d/%s", d.Host, d.Port, d.Name)
}

// 练习 25：结构体标签
// [TODO] APIConfig 带 json tag（字段名用 snake_case）
type APIConfig struct {
	// TODO: APIKey 导出为 "api_key", BaseURL 导出为 "base_url", Timeout 导出为 "timeout,omitempty"
	APIKey  string        `json:"api_key"`
	BaseURL string        `json:"base_url"`
	Timeout time.Duration `json:"timeout,omitempty"`
}

// ============================================================
// 练习 26：为 MemStore 添加 String() 方法
// 返回格式: "MemStore{keys: N}"
// 提示：在 methods.go 中给 *MemStore 添加 func (m *MemStore) String() string

func (m *MemStore) String() string {
	return fmt.Sprintf("MemStore{keys: %d}", len(m.data))
}

// ============================================================
// 练习 27：嵌套嵌入 + 方法提升
// Logger 嵌入 Name 字段和 io.Writer 行为
// ============================================================

type Prefix struct {
	Name string
}
type Logger struct {
	// TODO: 嵌入 Prefix 结构体（含 Name 字段），再添加 Writer io.Writer
	Prefix
	Writer io.Writer
}

// TODO: 实现 Log(msg string)，输出 "[Name] msg" 到 Writer
// 提示：fmt.Fprintf(logger.Writer, "[%s] %s\n", logger.Name, msg)

func (l *Logger) Log(msg string) {
	_, err := fmt.Fprintf(l.Writer, "[%s] %s\n", l.Name, msg)
	if err != nil {
		return
	}
}

// ============================================================
// 练习 28：结构体比较
// ============================================================

type Point struct{ X, Y int }

// TODO: 实现 Equal 方法，判断两个 Point 是否相同
func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// 思考题：如果 Point 里含有 map 或 slice 字段，还能用 == 比较吗？
// 答: 不能

// ============================================================
// 练习 29：工厂函数 + 校验
// NewServerConfig 返回 *ServerConfig，校验 host 非空、port > 0
// ============================================================

type ServerConfig struct {
	Host string
	Port int
}

// TODO: 校验 host != "" 且 port > 0，否则返回 error
func NewServerConfig(host string, port int) (*ServerConfig, error) {
	if host != "" && port > 0 {
		return &ServerConfig{host, port}, nil
	}
	return nil, fmt.Errorf("cannot create server config")
}

// ============================================================
// 练习 30：综合复习 — 接口 + 结构体 + 方法
// ============================================================

type Formatter interface {
	Format(v any) string
}

type JSONFormatter struct{}
type PrettyFormatter struct {
	Indent string
}

func (j JSONFormatter) Format(v any) string {
	formatted, err := json.Marshal(v)
	if err == nil {
		return string(formatted)
	}
	return ""
}

// TODO: PrettyFormatter 用 fmt.Sprintf 带缩进格式化
// 输出示例: "{\n  \"key\": \"value\"\n}"
func (p PrettyFormatter) Format(v any) string {
	formatted, err := json.MarshalIndent(v, "", p.Indent)
	if err == nil {
		return string(formatted)
	}
	return ""
}

// TODO: FormatAny 接收 any 和 Formatter，返回格式化结果
func FormatAny(v any, f Formatter) string {
	return f.Format(v)
}
