package practice

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// ============================================================
// net/http 练习（学完 review/http.go 后做）
// ============================================================

// 练习 42：HTTP GET 请求，返回 body
// [TODO] HTTPGet 发起 GET 请求，返回响应的 body 字符串和 error
// 提示: http.Get, io.ReadAll, resp.Body.Close()
func HTTPGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	rasp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rasp.Body)
	// 关于 Body.Close() 的错误：不需要处理。
	// 它的 error 只是底层连接关闭的状态，不影响你读到的数据。
	//官方推荐就当没返回值：defer resp.Body.Close()。

	buffer, err := io.ReadAll(rasp.Body)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

// 练习 43：简易 HTTP 服务端
// [TODO] StartServer 在指定端口启动服务，返回 /health 的 "OK" 响应
// 提示: http.HandleFunc, http.ListenAndServe
func StartServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// [TODO] 写入 "OK"
		_, err := fmt.Fprintf(w, "ok")
		if err != nil {
			return
		}
	})
	return http.ListenAndServe(":"+port, mux)
}

// 练习 44：带超时的 HTTP GET
// [TODO] HTTPGetWithTimeout 和 HTTPGet 一样，但超时后返回 timeout error
// 提示: http.Client{Timeout: ...}
func HTTPGetWithTimeout(url string, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	rasp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rasp.Body.Close()
	buffer, err := io.ReadAll(rasp.Body)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}
