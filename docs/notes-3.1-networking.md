# 3.1 网络编程学习笔记

> 本文件由 README.md 拆分而来（2026-09-07）。

### 🏁 三阶段 3.1 练习 (N-V)

- [x] N: TCP Echo Server
- [x] O: SimpleRouter — HTTP 路由
- [x] P: JSON API Server
- [x] Q: HTTP GET + 超时重试
- [x] R: HTTP 优雅关闭
- [x] S: 请求体限流中间件
- [x] T: SSE 流式响应
- [x] U: TCP 连接池
- [x] U-1: TCP 连接池（channel 版）
- [x] V: 文件上传

> 知识演示：`go run cmd/part3/main.go`

### 💡 3.1 关键结论（2026-08-16 沉淀）

- 连接池 `capacity` 是「最大空闲连接数」，不是总连接数：总连接 = 空闲 + 已借出，可超过 capacity
- `Close` 要幂等：先判断 `isClosed` 再 `close(ch)`，否则二次 Close 会 panic
- channel 版连接池：`Get`/`Put` 用 `select` + `default` 非阻塞取还，池满则关闭
- `net.Dial` 超时用 `net.Dialer{Timeout}` + `DialContext`，比 `DialTimeout` 更能响应 context 取消
- HTTP 响应一旦被 `Write`/`Flush` 提交，后续 `http.Error` 改状态码无效；SSE 事件分隔符是两个真实换行 `\n\n`

### 📚 3.1 补充学习（建议插入）

- [x] errgroup（`golang.org/x/sync/errgroup`）— 收集 goroutine 错误 + 协调取消
- [x] sync.WaitGroup — 等待一组 goroutine 全部退出
- [x] context 取消传播 — ctx.Done () / ctx.Err ()
- [x] net.ErrClosed + errors.Is — 区分正常关闭与真实错误
- [x] sync.Once — 保证关闭/初始化只执行一次
- [x] golang.org/x/net/netutil.LimitListener — 限制并发连接数

> 说明：以上 6 项已全部完成；net.ErrClosed / errgroup / LimitListener 在 practice/net_extra.go 中实现并验证。
