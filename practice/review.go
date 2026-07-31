package practice

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// ============================================================
// 复习题 1：完整 KVStore 测试套件
// 文件: kvstore_test.go 中追加 TestFileStore
// ============================================================
// [TODO] 参考 TestKVStore 的表驱动测试，为 FileStore 写同样三个测试用例
// 额外验证：关闭程序后重新 NewFileStore，数据应该还在（JSON 持久化验证）
// 提示: 用 os.CreateTemp 创建临时文件

// ============================================================
// 复习题 2：管道处理 — 读文件 → Plugin 转换 → 写文件
// ============================================================

// [TODO] ProcessFile 从 src 读取，依次用 plugins 处理每行，写入 dst
// 提示: os.ReadFile + strings.Split("\n") + 每行通过 RunPipeline + strings.Join + os.WriteFile
func ProcessFile(src, dst string, plugins []Plugin) error {
	srcTxt, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	strSlice := strings.Split(string(srcTxt), "\n")
	for index := range strSlice {
		strSlice[index] = RunPipeline(strSlice[index], plugins)
	}
	result := strings.Join(strSlice, "\n")
	err = os.WriteFile(dst, []byte(result), 0644) // 写入
	if err != nil {
		return err
	}
	return nil
}

// ============================================================
// 复习题 3：ConfigLoader — 反射 + JSON + default tag
// ============================================================
// [TODO] LoadConfig 读取 JSON 配置文件，反序列化到结构体，然后用反射的 FillDefaults 填充默认值
// 前提: 有 type ServerCfg struct { ... } 带 `default` tag
// 以下是一个参考结构体：

type ServerCfg struct {
	Host string `json:"host" default:"localhost"`
	Port int    `json:"port"`
}

// [TODO] LoadServerConfig 从 path 读取 JSON 到 ServerCfg，调用 FillDefaults 补全默认值
// 返回 *ServerCfg 和 error
func LoadServerConfig(path string) (*ServerCfg, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &ServerCfg{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	FillDefaults(cfg)
	return cfg, nil

}

// ============================================================
// 复习题 A：健康检查器 — 定时 GET /health，超时报警
// ============================================================

// [TODO] HealthCheck 每隔 interval 请求 url，连续 failCount 次失败返回 error
// 每次请求用 timeout 作超时限制，成功打印 "OK"，失败打印 "FAIL"
// 提示: time.NewTicker + HTTPGetWithTimeout
func HealthCheck(url string, interval, timeout time.Duration, failCount int) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	client := &http.Client{
		Timeout: timeout,
	}
	actualFailCount := 0
	for range ticker.C {
		if actualFailCount >= failCount {
			break
		}
		r, err2 := client.Get(url)
		if err2 != nil {
			fmt.Println("FAIL")
			actualFailCount++
			continue
		}
		r.Body.Close()
		fmt.Println("ok")
		actualFailCount = 0
	}
	return errors.New("health check fail")
}

// ============================================================
// 复习题 B：JSON API 客户端
// ============================================================

// [TODO] FetchAndParse GET 请求 JSON API，反序列化到结构体，填默认值
// 示例结构体已经定义（ServerCfg），你可以直接用
// 提示: HTTPGet + json.Unmarshal + FillDefaults(ptr)
func FetchAndParse(url string, timeout time.Duration, target interface{}) error {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return err
	}
	FillDefaults(target)
	return nil
}

// ============================================================
// 复习题 C：数据管道 — Reader → 过滤 → 转换 → Writer
// ============================================================

// [TODO] TransformStream 从 r 逐行读取，保留含 keyword 的行，每行用 plugins 处理，写入 w
// 组合了: bufio.Scanner + FilterLines(逻辑复用) + RunPipeline + fmt.Fprintf
// 提示：直接在一个循环里做：Scan → Contains → RunPipeline → Fprintf
func TransformStream(r io.Reader, w io.Writer, keyword string, plugins []Plugin) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, keyword) {
			_, err := fmt.Fprintf(w, "%s\n", RunPipeline(line, plugins))
			if err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

// ============================================================
// 二阶段复习题（进入第三阶段前磨刀）
// ============================================================

// 复习题 D：结构体组合 + 方法重写
// ============================================================
// [TODO] TimestampLogger 嵌入 io.Writer，Write 方法改成带时间戳的格式
// Write 应输出 "[2006-01-02 15:04:05] <data>" 到嵌入的 Writer
type TimestampLogger struct {
	io.Writer
}

func NewTimestampLogger(w io.Writer) *TimestampLogger {
	return &TimestampLogger{Writer: w}
}

func (l *TimestampLogger) Write(data []byte) (int, error) {
	_, err := l.Writer.Write(
		[]byte(
			fmt.Sprintf("%s %s", time.Now().Format("[2006-01-02 15:04:05] "), string(data))))
	if err == nil {
		return 0, err
	}
	return 0, errors.New("not implemented")
}

// 复习题 E：类型开关 — 格式化打印不同类型
// ============================================================
// [TODO] Describe 用 type switch 判断 v 的类型，返回描述字符串
// string → "string: <value>"
// int → "int: <value>"
// []int → "slice: len=N"
// 其他 → "unknown"
func Describe(v interface{}) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("string: %v", v.(string))
	case int:
		return fmt.Sprintf("int: %v", v.(int))
	case []int:
		return fmt.Sprintf("slice: len=%d", len(v.([]int)))
	default:
		return fmt.Sprintf("unknown")
	}
}

// 复习题 F：反射 — 按名称设置结构体字段
// ============================================================
// [TODO] SetField 将结构体指针 ptr 中名为 name 的字段设为 value（字符串）
// 前提: ptr 必须是 *struct，且该字段必须是 string 类型
// 提示: reflect.ValueOf(ptr).Elem().FieldByName(name).SetString(value)
func SetField(ptr interface{}, name, value string) error {
	//val := reflect.ValueOf(ptr).Elem()
	//typ := val.Type()
	//for i := 0; i < typ.NumField(); i++ {
	//	vField := val.Field(i)
	//	tField := typ.Field(i)
	//	if tField.Name == name && vField.Kind() == reflect.String && vField.CanSet() {
	//		vField.SetString(value)
	//	}
	//}
	//return nil
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return errors.New("not a pointer or struct")
	}
	field := val.Elem().FieldByName(name)
	if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.String {
		return errors.New("not a valid field")
	}
	field.SetString(value)
	return nil
}

// 复习题 G：实现 io.Writer 统计写入量
// ============================================================
// [TODO] CountingWriter 包装一个 io.Writer，统计写入的总字节数
// Total() 返回总字节数
type CountingWriter struct {
	w     io.Writer
	total int64
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{w: w, total: 0}
}

func (cw *CountingWriter) Write(p []byte) (int, error) {
	count, err := cw.w.Write(p)
	cw.total += int64(count)
	return count, err

}

func (cw *CountingWriter) Total() int64 {
	return cw.total
}

// 复习题 H：HTTP 中间件链
// ============================================================
// [TODO] ChainMiddleware 将多个中间件串成链
// 调用顺序: mw1 → mw2 → handler
// 提示: 每个中间件签名为 func(http.Handler) http.Handler，从末尾反向包装
type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// 复习题 I：防止 goroutine 泄漏
// ============================================================
// [TODO] Search 并发查询多个数据源，返回第一个成功的结果，并取消其他查询
// 用 done channel 模式: 每个查询 goroutine 先检查 done 是否关闭再干活
// 提示: 用 buffered chan (result + done) 经典模式，避免泄漏
type SearchResult struct {
	Data string
	Err  error
}

func Search(sources []func() (string, error)) (string, error) {
	resultChan := make(chan SearchResult, 1)
	doneChan := make(chan bool)
	var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)
		go func(source func() (string, error)) {
			defer wg.Done()
			data, err := source()
			if err != nil {
				return
			}
			sResult := SearchResult{data, err}
			select {
			case <-doneChan:
				return
			case resultChan <- sResult:
				return
			}
		}(source)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	if final, ok := <-resultChan; !ok {
		close(doneChan)
		return "", nil
	} else {
		close(doneChan)
		return final.Data, final.Err
	}

}

// 复习题 J：select 合并多个 channel（Fan-in）
// ============================================================
// [TODO] FanIn 用 select 将多个 channel 合并为一个
// 输入 channel 可能随时关闭，输出 channel 在所有输入关闭后关闭
func FanIn(chs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// 复习题 K：Context 树 — 父 cancel 传播到子
// ============================================================
// [TODO] WithChild 创建父子两个 context.CancelFunc
// 父 cancel 会同时取消子（Go 已保证），但子 cancel 不应该影响父
// 返回 parentCancel 和 childCancel，以及子 ctx
// 提示: 直接用 context.WithCancel
func WithChild() (parentCtx context.Context, parentCancel context.CancelFunc, childCtx context.Context, childCancel context.CancelFunc) {
	parentCtx, parentCancel = context.WithCancel(context.Background()) // empty ctx
	childCtx, childCancel = context.WithCancel(parentCtx)
	return
}

// 复习题 L：错误包装与解包
// ============================================================
// [TODO] WrapError 用 fmt.Errorf("%w") 包装错误
// UnwrapTo 用 errors.As 检查 err 链中是否有 *target 类型的错误
var ErrNotFound = errors.New("not found")
var ErrTimeout = errors.New("timeout")

func WrapError(base error, msg string) error {
	return fmt.Errorf("%s: %w", msg, base)
}

func UnwrapTo(err error, target interface{}) bool {
	return errors.As(err, &target)
}

// 复习题 M：带 Context 取消的 Pipeline
// ============================================================
// [TODO] PipelineWithCtx
// stage1: 接收 nums，遇到 ctx.Done() 时停止发送
// stage2: 从 stage1 读并平方，遇到 ctx.Done() 时停止
// 返回 stage2 的 channel
// 提示: for { select { case <-ctx.Done(): return; case ch <- val: ... } }
func PipelineWithCtx(ctx context.Context, nums []int) <-chan int {
	inChan := make(chan int)
	go func(ctx context.Context, nums []int) {
		defer close(inChan)
		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case inChan <- num:
			}
		}
	}(ctx, nums)
	outChan := make(chan int)
	go func(in <-chan int) {
		defer close(outChan)
		for {
			select {
			case <-ctx.Done():
				return
			case num, ok := <-in:
				if !ok {
					return
				}
				outChan <- num * num
			}
		}
	}(inChan)
	return outChan
}
