package practice

import (
	"errors"
	"io"
	"net/http"
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
	return errors.New("not implemented")
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
	return nil
}

// 练习 Q：HTTP 客户端超时 + 重试
// ============================================================
// [TODO] HTTPGetWithRetry GET 请求 url，最多重试 retries 次
// 每次请求有 timeout 超时，非 2xx 或网络错误都算失败
// 全部失败返回最后一次的错误
// 提示: http.Client{Timeout: timeout}.Do 或 Get
func HTTPGetWithRetry(url string, timeoutMs int, retries int) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

// 练习 R：HTTP Server 优雅关闭
// ============================================================
// [TODO] GracefulServer 启动 HTTP 服务，返回一个 stop 函数
// stop 调用后服务应优雅关闭（不再接受新请求，等待正在处理的请求完成）
// 提示: http.Server{Addr: addr}, srv.ListenAndServe(), srv.Shutdown(context.Background())
func GracefulServer(addr string, handler http.Handler) (stop func() error, err error) {
	return nil, errors.New("not implemented")
}

// 练习 S：HTTP 请求体限流（MiddleWare）
// ============================================================
// [TODO] LimitBodyMiddleware 返回一个中间件，限制请求体大小不超过 maxBytes
// 超过则返回 413 Request Entity Too Large
// 提示: http.MaxBytesReader
func LimitBodyMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return nil
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
}

func NewConnPool(addr string, capacity int) *ConnPool {
	return nil
}

func (p *ConnPool) Get() (io.ReadWriteCloser, error) {
	return nil, errors.New("not implemented")
}

func (p *ConnPool) Put(conn io.ReadWriteCloser) {
}

func (p *ConnPool) Close() {
}

// 练习 V：HTTP 文件上传服务
// ============================================================
// [TODO] UploadHandler 处理 multipart/form-data 文件上传
// 读取 "file" 字段的文件，保存到 uploadDir/原文件名
// 返回 JSON: {"filename": "...", "size": N}
// 提示: req.ParseMultipartForm(maxMemory), req.FormFile("file"), io.Copy
func UploadHandler(uploadDir string) http.HandlerFunc {
	return nil
}
