package review

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ============================================================
// 补充模块：net/http 基础
// ============================================================

func HTTPGetDemo() {
	fmt.Println("=== GET 请求 ===")

	// 最简单的 GET
	resp, err := http.Get("https://httpbin.org/get?name=gopher")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("Body 前 100 字节: %s...\n", body[:min(len(body), 100)])
}

func HTTPPostDemo() {
	fmt.Println("\n=== POST 请求 ===")

	// POST JSON
	jsonBody := `{"name":"gopher","age":3}`
	resp, err := http.Post(
		"https://httpbin.org/post",
		"application/json",
		strings.NewReader(jsonBody),
	)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("返回: %s...\n", body[:min(len(body), 100)])
}

func HTTPNewRequestDemo() {
	fmt.Println("\n=== http.NewRequest（自定义 Header）===")

	// NewRequest 可以加 Header、改 Method
	req, _ := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	req.Header.Set("X-Custom", "learngo")
	req.Header.Set("Authorization", "Bearer token123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("自定义 Header 回声: %s", body)
}

func HTTPServerDemo() {
	fmt.Println("\n=== HTTP 服务端 ===")

	// 最简服务端：一个路径，一个处理函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s! Method=%s Body=%s", r.URL.Path, r.Method, r.Body)
	}

	http.HandleFunc("/hello", handler)

	// 启动服务（在 goroutine 里，避免阻塞演示）
	go func() {
		if err := http.ListenAndServe(":9999", nil); err != nil {
			fmt.Println("服务端错误:", err)
		}
	}()

	// 给服务端一点时间启动
	time.Sleep(100 * time.Millisecond)

	// 客户端请求验证
	resp, _ := http.Get("http://localhost:9999/hello")
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("服务端响应: %s\n", body)
}

func HTTPTimeoutDemo() {
	fmt.Println("\n=== HTTP 超时设置 ===")

	client := &http.Client{
		Timeout: 2 * time.Second, // 整个请求的超时（包含 DNS + 连接 + 读取）
	}

	start := time.Now()
	_, err := client.Get("https://httpbin.org/delay/1") // 1 秒后返回
	fmt.Printf("延迟 1s → 耗时 %v, err=%v\n", time.Since(start), err)

	start = time.Now()
	_, err = client.Get("https://httpbin.org/delay/5") // 5 秒，超过 2s 超时
	fmt.Printf("延迟 5s → 耗时 %v, err=%v\n", time.Since(start), err)
}
