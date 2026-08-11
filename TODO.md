# Go 学习路线图

> 基于 [《The Way To Go》](https://learnku.com/docs/the-way-to-go) 教程，目标：为学习**分布式系统**打下坚实基础。

---

> ⚠️ **AI 助手指令（每次会话必读）**
>
> 这是一个 Go 学习项目，用户正在按路线图自学。你的角色是**导师**，不是代劳者。
>
> - **禁止直接实现 TODO/练习**：不要替用户写代码，让用户自己动手。
> - **可以做的事**：解释概念、分析已有代码、回答疑问、指出错误、给出思路提示、运行验证命令。
> - **如果用户卡住**：先给思路提示，不要直接给完整答案。用户要求看答案时再展示。
> - **适当扩展学习**：发现用户缺少前置知识时，主动建议插入补充模块。
> - **检查练习**：按用户要求检查相关习题。
> - **关键文件**：
> >
>   - `review/basics.go` … `review/http.go` — 一、二阶段知识点讲解
>   - `review/networking.go` — 三阶段网络编程知识点讲解
>   - `practice/exercises.go` — 练习 1-9（基础）
>   - `practice/interfaces.go` — 练习 15-22（接口）
>   - `practice/io.go` — 练习 23（文件/JSON）
>   - `practice/structs.go` — 练习 24-30（结构体）
>   - `practice/reflect.go` — 练习 31-33, 36（反射 + panic）
>   - `practice/kvstore_test.go` — 练习 34-35（测试）
>   - `practice/time.go` — 练习 37-38（time 包）
>   - `practice/io_reader.go` — 练习 39-41（io.Reader/Writer）
>   - `practice/http.go` — 练习 42-44（HTTP）
>   - `practice/goroutine.go` — 练习 45-46（goroutine/channel）
>   - `practice/concurrency.go` — 练习 10-14, 47-53（并发综合）
>   - `practice/review.go` — 复习题 A-M（综合）
>   - `practice/networking.go` — 练习 N-V（网络编程）
>   - `cmd/part1/` / `cmd/part2/` / `cmd/part3/` — 按阶段运行演示
> > - **运行方式**：
      > >
- `go run cmd/part3/main.go` — 运行第三阶段知识演示
> >   - 练习验证：直接 `go build ./practice/...` 检查编译

---

## 🎯 当前进度

| 阶段       | 状态               |
|----------|------------------|
| 第一阶段     | ✅                |
| 第二阶段     | ✅                |
| **第三阶段** | **🔄 3.1 网络编程中** |
| 第四阶段     | 待开始              |

> **上次会话：2026-07-31** | 二阶段复习完成，开始三阶段网络编程

---

### 🏁 三阶段 3.1 练习 (N-V)

- [ ] N: TCP Echo Server
- [ ] O: SimpleRouter — HTTP 路由
- [ ] P: JSON API Server
- [ ] Q: HTTP GET + 超时重试
- [ ] R: HTTP 优雅关闭
- [ ] S: 请求体限流中间件
- [ ] T: SSE 流式响应
- [ ] U: TCP 连接池
- [ ] V: 文件上传

> 知识演示：`go run cmd/part3/main.go`

### 📚 3.1 补充学习（建议插入）

- [ ] errgroup（`golang.org/x/sync/errgroup`）— 收集 goroutine 错误 + 协调取消
- [ ] sync.WaitGroup — 等待一组 goroutine 全部退出
- [ ] context 取消传播 — ctx.Done() / ctx.Err()
- [ ] net.ErrClosed + errors.Is — 区分正常关闭与真实错误
- [ ] sync.Once — 保证关闭/初始化只执行一次
- [ ] golang.org/x/net/netutil.LimitListener — 限制并发连接数

---

## 第三阶段：分布式系统专项 🎯

- [ ] 3.1 网络编程（TCP/HTTP Server/Client）🔄
- [ ] 3.2 gRPC（protobuf + Unary/Stream RPC）
- [ ] 3.3 服务注册与发现（etcd）
- [ ] 3.4 负载均衡
- [ ] 3.5 分布式共识（Raft）
- [ ] 3.6 消息队列（Pub/Sub）
- [ ] 3.7 可观测性（OpenTelemetry）

## 第四阶段：综合项目 🏗️

- [ ] 4.1 分布式 KV 存储
- [ ] 4.2 任务调度系统

## 推荐资源

- **并发**: [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- **Raft**: [Raft 可视化](https://raft.github.io/)
- **分布式**: [MIT 6.824](https://pdos.csail.mit.edu/6.824/)
- **gRPC**: [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
