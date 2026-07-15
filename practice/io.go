package practice

import (
	"encoding/json"
	"os"
)

// ============================================================
// 文件与 JSON 练习（学完 review/io_json.go 后做）
// ============================================================

// 练习 23：JSON 读写
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

func SaveConfig(path string, cfg AppConfig) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func LoadConfig(path string) (AppConfig, error) {
	var cfg AppConfig
	b, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return AppConfig{}, err
	}
	return cfg, nil
}
