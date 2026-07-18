package practice

import (
	"encoding/json"
	"os"
	"strings"
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
