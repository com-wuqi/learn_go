package practice

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
