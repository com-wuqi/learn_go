package review

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"
)

// ============================================================
// 第三阶段 3.1：网络编程 — 知识讲解
// ============================================================

// --- TCP 基础：最简 Echo 服务端 ---
// net.Listen("tcp", addr)  → 返回 Listener
// listener.Accept()        → 阻塞等待连接，返回 net.Conn
// conn.Read(buf)           → 从连接读数据
// conn.Write(data)         → 向连接写数据
// io.Copy(dst, src)        → 最简 echo: io.Copy(conn, conn)

func TCPEchoDemo() {
	fmt.Println("=== TCP Echo 服务端 ===")

	// 1. 监听端口
	listener, err := net.Listen("tcp", ":0") // :0 表示随机端口
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("  监听: %s\n", listener.Addr())

	go func() {
		for {
			// 2. 循环 Accept
			conn, err := listener.Accept()
			if err != nil {
				return // 通常是 listener 关闭
			}
			// 3. 每个连接一个 goroutine
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c) // 读什么写什么
			}(conn)
		}
	}()

	// 客户端验证
	conn, _ := net.Dial("tcp", listener.Addr().String())
	fmt.Fprintf(conn, "hello tcp\n")
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	fmt.Printf("  client 收到: %s", resp)
}

// --- TCP 粘包/半包问题 ---
// TCP 是流式协议，没有消息边界。Read 不保证一次读完一条"消息"。
// 解决方案: 固定长度 / 分隔符 / 长度前缀。
// 这里演示分隔符方案：\n 作为消息边界。

func TCPFrameDemo() {
	fmt.Println("\n=== TCP 消息边界（分隔符方案）===")

	listener, _ := net.Listen("tcp", ":0") // 监听本机所有接口上的一个随机空闲端口。
	fmt.Printf("addr is: %s\n", listener.Addr().String())
	defer listener.Close()
	addr := listener.Addr().String()

	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n') // 读到 \n 才返回
			if err != nil {
				return
			}
			fmt.Fprintf(conn, "echo: %s", line)
		}
	}()

	conn, _ := net.Dial("tcp", addr)
	fmt.Fprintf(conn, "msg1\nmsg2\n") // 两条消息一起发
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	resp2, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	fmt.Printf("  收到: %q %q\n", resp, resp2) // echo: msg1 echo: msg2
}

// --- http.Handler 接口 ---
// 所有 HTTP 服务都是围绕这个接口：
//
//   type Handler interface {
//       ServeHTTP(ResponseWriter, *Request)
//   }
//
// http.HandleFunc("/path", fn) 内部是把 fn 转成 HandlerFunc 类型，
// 而 HandlerFunc 实现了 Handler 接口（ServeHTTP 调自己）。
//
// 自定义 Router 就是实现了 Handler 接口的结构体。

func HandlerInterfaceDemo() {
	fmt.Println("\n=== http.Handler 接口 ===")

	// 方式1: 闭包 → http.HandlerFunc 自动适配
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	// 方式2: http.HandlerFunc 是 func(http.ResponseWriter, *http.Request)  的适配器
	// 它实现了 http.Handler 接口：ServeHTTP 方法调用函数自身
	mux := http.NewServeMux()
	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hi")
	})

	go http.ListenAndServe(":0", mux)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("  HandlerFunc 实现了 Handler 接口（ServeHTTP 调用函数自身）")
}

// --- http.Server 优雅关闭 ---
// srv.Shutdown(ctx)  vs  srv.Close()
// Shutdown: 不再接受新连接，等现有请求处理完
// Close:    立即关闭所有连接

func GracefulShutdownDemo() {
	fmt.Println("\n=== 优雅关闭 ===")

	srv := &http.Server{Addr: ":0"}
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "done")
	})
	http.HandleFunc("/fast", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "fast")
	})

	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	// 模拟优雅关闭
	go func() {
		time.Sleep(50 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("  调用 Shutdown...")
		srv.Shutdown(ctx)
		fmt.Println("  Shutdown 完成")
	}()

	time.Sleep(300 * time.Millisecond)
	fmt.Println("  要点: Shutdown 会等待处理中的请求完成")
}

// --- http.Flusher：流式响应 / SSE ---
// ResponseWriter 默认会缓冲。要实现实时推送，需要 Flusher：
//
//   flusher, ok := w.(http.Flusher)
//   fmt.Fprintf(w, "data: ...\n\n")
//   flusher.Flush()  // 立即发送到客户端

func FlusherDemo() {
	fmt.Println("\n=== http.Flusher（SSE 风格）===")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		flusher, ok := w.(http.Flusher)
		//flusher, ok := http.NewResponseController(w).(http.Flusher)

		if !ok {
			http.Error(w, "streaming not supported", 500)
			return
		}

		for i := 1; i <= 3; i++ {
			select {
			case <-r.Context().Done():
				return // 客户端断开
			default:
				fmt.Fprintf(w, "data: tick %d\n\n", i)
				flusher.Flush()
				time.Sleep(300 * time.Millisecond)
			}
		}
	} // codex resume 019fd4dd-ee09-71b1-b62d-a35074646dfa

	go http.ListenAndServe(":0", http.HandlerFunc(handler))
	time.Sleep(100 * time.Millisecond)
	fmt.Println("  要点: Flusher.Flush() 立即推送 / Context().Done() 检测客户端断开")
}

// --- http.MaxBytesReader：限制请求体 ---
func MaxBytesDemo() {
	fmt.Println("\n=== 请求体大小限制 ===")

	handler := func(w http.ResponseWriter, r *http.Request) {
		// 限制 Body 最多 100 字节
		r.Body = http.MaxBytesReader(w, r.Body, 100)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			// MaxBytesReader 会在超限时返回错误
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}
		fmt.Fprintf(w, "received %d bytes", len(data))
	}

	go http.ListenAndServe(":0", http.HandlerFunc(handler))
	time.Sleep(100 * time.Millisecond)
	fmt.Println("  要点: MaxBytesReader 包装 Body，超限自动 413")
}

// --- multipart/form-data 文件上传 ---
// 客户端: multipart.Writer 生成 multipart body
//
//	每个 part 带 Content-Disposition: form-data; name="file"; filename="xxx"
//
// 服务端: r.ParseMultipartForm(maxMemory) 解析 → r.FormFile("file") 取出文件
//
//	注意: FormFile 返回的 file 必须 Close()；原始文件名在 FileHeader 里
func UploadDemo() {
	fmt.Println("\n=== 文件上传 (multipart/form-data) ===")

	// ---------- 服务端 ----------
	uploadHandler := func(w http.ResponseWriter, r *http.Request) {
		// 1. 先限制请求体大小（防滥用）
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

		// 2. 解析 multipart 表单
		if err := r.ParseMultipartForm(4 << 20); err != nil {
			// 区分超限与格式错误：前者回 413，后者回 400
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "parse form failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 3. 取出文件（FormFile 内部 = 第一个 name 为 "file" 的文件 part）
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "no file field: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 4. 落盘保存（演示用临时文件；真实项目应放到 uploads 目录并管理清理）
		dst, err := os.CreateTemp("", "upload-*")
		if err != nil {
			http.Error(w, "create file failed", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		n, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, "save file failed", http.StatusInternalServerError)
			return
		}

		// 5. 返回 JSON（原始文件名在 FileHeader 里）
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"filename": header.Filename,
			"size":     n,
			"saved":    dst.Name(),
		})
	}

	// 显式创建 listener，客户端才能拿到真实端口
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	srv := &http.Server{Handler: http.HandlerFunc(uploadHandler)}
	go srv.Serve(ln)
	defer srv.Close()
	time.Sleep(100 * time.Millisecond)

	// ---------- 客户端 ----------
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	part, _ := mw.CreateFormFile("file", "hello.txt") // 字段名必须和服务端 FormFile 一致
	part.Write([]byte("hello multipart!\n"))

	mw.Close() // 必须 Close，否则 body 缺收尾 boundary

	req, _ := http.NewRequest("POST", "http://"+ln.Addr().String()+"/upload", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("  请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("  客户端发送: hello.txt")
	fmt.Printf("  服务端返回: %s\n", body)

	// 超限验证：12MB > 10MB 上限，服务端应回 413
	var bigBuf bytes.Buffer
	bigMw := multipart.NewWriter(&bigBuf)
	bigPart, _ := bigMw.CreateFormFile("file", "big.txt")
	bigPart.Write(bytes.Repeat([]byte("x"), 12<<20))
	bigMw.Close()

	bigReq, _ := http.NewRequest("POST", "http://"+ln.Addr().String()+"/upload", &bigBuf)
	bigReq.Header.Set("Content-Type", bigMw.FormDataContentType())

	bigResp, err := http.DefaultClient.Do(bigReq)
	if err != nil {
		fmt.Println("  超限请求失败:", err)
		return
	}
	defer bigResp.Body.Close()
	bigBody, _ := io.ReadAll(bigResp.Body)
	fmt.Printf("  超限请求: status=%d body=%q\n", bigResp.StatusCode, bigBody)
}

// ============================================================
// RunPhase3 执行第三阶段演示
// ============================================================

func RunPhase3() {
	TCPEchoDemo()
	TCPFrameDemo()
	HandlerInterfaceDemo()
	GracefulShutdownDemo()
	FlusherDemo()
	MaxBytesDemo()
	UploadDemo()
}
