package practice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ============================================================
// 第三阶段 3.1：网络编程
// ============================================================

// 练习 N：TCP Echo 服务端
// ============================================================
// [TODO] StartEchoServer 在指定端口启动 TCP echo 服务
// 每收到一个连接，循环读取直到 EOF，把读到的内容原样写回
// 提示: net.Listen("tcp", addr) → 循环 Accept() → go handle(conn)
func StartEchoServer(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("Error closing listener")
		}
	}(listen)
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("Error accepting connection")
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				_, err := io.Copy(conn, conn)
				if err != nil {
					return
				}
			}(conn)
		}
	}()

	return nil
}

// 练习 O：简易 HTTP Router
// ============================================================
// [TODO] SimpleRouter 实现一个基于 map 的路由器
// 注册 route → handler，ServeHTTP 根据请求 path 分发
// 找不到路由返回 404
// 提示: 实现 http.Handler 接口，内部用 map[string]http.HandlerFunc
type SimpleRouter struct {
	routes map[string]http.HandlerFunc
}

func NewSimpleRouter() *SimpleRouter {
	return &SimpleRouter{routes: make(map[string]http.HandlerFunc)}
}

func (r *SimpleRouter) Handle(path string, handler http.HandlerFunc) {
	r.routes[path] = handler
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO
	if handler, ok := r.routes[req.URL.Path]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// 练习 P：JSON API 服务端
// ============================================================
// [TODO] JSONServer 用 SimpleRouter 搭建两个接口：
//
//	GET  /health → {"status": "ok"}
//	POST /echo   → 读取请求体 JSON，原样返回
//
// 返回 *SimpleRouter 实例
func NewJSONServer() *SimpleRouter {
	health := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	echo := func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		data, _ := io.ReadAll(req.Body)
		var encodedData any
		if err := json.Unmarshal(data, &encodedData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := json.NewEncoder(w).Encode(encodedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	return &SimpleRouter{routes: map[string]http.HandlerFunc{"/echo": echo, "/health": health}}
}

// 练习 Q：HTTP 客户端超时 + 重试
// ============================================================
// [TODO] HTTPGetWithRetry GET 请求 url，最多重试 retries 次
// 每次请求有 timeout 超时，非 2xx 或网络错误都算失败
// 全部失败返回最后一次的错误
// 提示: http.Client{Timeout: timeout}.Do 或 Get
func HTTPGetWithRetry(url string, timeoutMs int, retries int) (*http.Response, error) {
	// 最多尝试 retries+1 次
	client := http.Client{
		Timeout: time.Duration(time.Duration(timeoutMs) * time.Millisecond),
	}
	var lastErr error
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < retries+1; i++ {
		resp, err2 := client.Do(req)
		if err2 != nil {
			lastErr = err2
			continue
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
		resp.Body.Close()
		// lastErr = err2 // 此处非网络层故障，nil
		lastErr = errors.New(resp.Status)
	}
	return nil, lastErr
}

// 练习 R：HTTP Server 优雅关闭
// ============================================================
// [TODO] GracefulServer 启动 HTTP 服务，返回一个 stop 函数
// stop 调用后服务应优雅关闭（不再接受新请求，等待正在处理的请求完成）
// 提示: http.Server{Addr: addr}, srv.ListenAndServe(), srv.Shutdown(context.Background())
func GracefulServer(addr string, handler http.Handler) (stop func() error, err error) {
	server := &http.Server{Handler: handler}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go server.Serve(ln)
	stop = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err1 := server.Shutdown(ctx); err1 != nil {
			server.Close()
			return err1
		}
		return nil
	}
	return stop, nil
}

// 练习 S：HTTP 请求体限流（MiddleWare）
// ============================================================
// [TODO] LimitBodyMiddleware 返回一个中间件，限制请求体大小不超过 maxBytes
// 超过则返回 413 Request Entity Too Large
// 提示: http.MaxBytesReader
func LimitBodyMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	limitFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//defer r.Body.Close()
			if r.ContentLength > maxBytes {
				http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
				return // 必须，否则会污染相映
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
	return limitFunc
}

// 练习 T：流式 JSON 响应（Server-Sent Events 风格）
// ============================================================
// [TODO] StreamHandler 持续向客户端发送 "data: tick\n\n" 每秒一次
// 客户端断开连接时停止（ctx.Done()）
// 提示: w.Header().Set("Content-Type", "text/event-stream")
//
//	用 flusher, ok := w.(http.Flusher); ok { flusher.Flush() }
//	监听 req.Context().Done()
func StreamHandler(w http.ResponseWriter, req *http.Request) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	//ctx, cancel := context.WithCancel(req.Context())
	//defer cancel() // cancel 无效
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache") // 避免中间层缓存
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case <-ticker.C:
			_, err := io.WriteString(w, "data: tick\n\n")
			if err != nil {
				//http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			flusher.Flush()
		case <-req.Context().Done():
			return
		}
	}

}

// 练习 U：连接池 TCP Client
// ============================================================
// [TODO] ConnPool 管理一组可复用的 TCP 连接
// Get() 获取一个连接（无可用时创建新的）
// Put(conn) 归还连接（池满则关闭）
// Close() 关闭所有连接
// 提示: 用 chan net.Conn 作缓冲池
type ConnPool struct {
	// TODO
	addr     string
	capacity int
	conns    []net.Conn
	lock     sync.Mutex
}

func NewConnPool(addr string, capacity int) *ConnPool {
	conns := make([]net.Conn, 0, capacity)
	var isOk bool
	isOk = true
	for i := 0; i < capacity; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			isOk = false
			break
		}
		conns = append(conns, conn)
	}
	if !isOk {
		for _, conn := range conns {
			_ = conn.Close()
		}
		return nil
	}
	return &ConnPool{addr: addr, conns: conns, capacity: capacity}
}

func (p *ConnPool) Get() (net.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.conns) == 0 {
		conn, err := net.Dial("tcp", p.addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	conn := p.conns[0]
	p.conns = p.conns[1:]
	return conn, nil
}

func (p *ConnPool) Put(conn net.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.conns) == p.capacity {
		_ = conn.Close()
		return
	}
	p.conns = append(p.conns, conn)

}

func (p *ConnPool) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, conn := range p.conns {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection")
			continue
		}
	}
}

// 练习 U-1：连接池 TCP Client（channel 版）
// ============================================================
// [TODO] 用 chan net.Conn 作缓冲池，实现懒建版连接池
// NewConnPoolChan 只创建缓冲 channel，不拨号
// Get() 无空闲连接时 net.Dial 新建
// Put(conn) 归还连接，池满则关闭
// Close() 关闭所有空闲连接（注意 close(ch) 与 Put 的并发）
type ConnPoolChan struct {
	// TODO
	addr     string
	capacity int
	conns    chan net.Conn
	lock     sync.Mutex
	isClosed bool
}

func NewConnPoolChan(addr string, capacity int) *ConnPoolChan {
	return &ConnPoolChan{addr: addr, capacity: capacity, conns: make(chan net.Conn, capacity), isClosed: false, lock: sync.Mutex{}}
}

func (p *ConnPoolChan) Get() (net.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.isClosed {
		return nil, errors.New("connection pool is closed")
	}
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		conn, err := net.Dial("tcp", p.addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (p *ConnPoolChan) Put(conn net.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.isClosed {
		_ = conn.Close()
		return
	}
	select {
	case p.conns <- conn:
	default:
		_ = conn.Close()
	}
}

func (p *ConnPoolChan) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.isClosed {
		return // Close 的幂等
	}
	p.isClosed = true
	close(p.conns)
	for conn := range p.conns {
		_ = conn.Close()
	}
}

// 练习 V：HTTP 文件上传服务
// ============================================================
// [TODO] UploadHandler 处理 multipart/form-data 文件上传
// 读取 "file" 字段的文件，保存到 uploadDir/原文件名
// 返回 JSON: {"filename": "...", "size": N}
// 提示: req.ParseMultipartForm(maxMemory), req.FormFile("file"), io.Copy
func UploadHandler(uploadDir string) http.HandlerFunc {
	uploadHandler := func(w http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(w, req.Body, 10<<20)
		if err := req.ParseMultipartForm(2 << 20); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(w, maxErr.Error(), http.StatusRequestEntityTooLarge)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		file, header, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Error closing file")
			}
		}(file)
		ext := filepath.Ext(header.Filename) // 后缀
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst, err := os.Create(filepath.Join(uploadDir, name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close() // 关闭
		n, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"filename":     name,
			"originalName": header.Filename,
			"size":         n,
		})
		if err != nil {
			fmt.Println(err)
		}
		return

	}
	return uploadHandler
}
