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
> - **当前进度快照（2026-08-17）**：
> >
>   - 第三阶段 3.1 网络编程与补充学习已完成 ✅；3.2 gRPC 工具链已就绪，3.2-A 演示完成、练习进行中
>   - 补充学习已全部完成：errgroup / WaitGroup / context 取消 / net.ErrClosed / sync.Once / LimitListener
> - **项目约定（用户 2026-08-11 明确）**：
> >
>   - 练习题解答写在 `practice/*.go`（用户自己动手）；`review/*.go` 只做大范围复习讲解
>   - 检查练习：先读代码 → `go build ./...` → 必要时写临时测试（`practice/zz_check_test.go`，测完删除）
> - **环境/工具备忘**：
> >
>   - 沙箱 bwrap 故障（2026-08-17 实测）：`apply_patch` 新增文件可用，但修改/删除已有文件会报 bwrap 错误；用提权
      `exec_command` + Python 定点替换文件，改完 `gofmt -w` 并编译验证
>   - gRPC 工具链（用户本地 /home/composer/go/bin，不在默认 PATH）：protoc 35.1、protoc-gen-go 1.36.12、protoc-gen-go-grpc
      1.6.2；使用前 export PATH 或全路径
>   - 环境有 `HTTP_PROXY=http://127.0.0.1:7897/` 且 `NO_PROXY` 不含 `[::]`：本地 HTTP 演示客户端地址用 `127.0.0.1`（绑
      `:0` 会被代理返回 502）
>   - `review/networking.go` FlusherDemo 有一行 `} // codex resume <uuid>` 残留注释（未确认是否删除）；UploadDemo 的
      `ParseMultipartForm(4 << 20)` 是用户有意改的（演示大文件落盘）
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
>   - `practice/networking.go` — 练习 N-V + U-1（网络编程）
>   - `practice/net_extra.go` — 3.1 补充练习（net.ErrClosed / errgroup / LimitListener）
>   - `review/net_extra.go` — 3.1 补充知识讲解（net.ErrClosed / errgroup / LimitListener）
>   - `cmd/part1/` / `cmd/part2/` / `cmd/part3/` — 按阶段运行演示
> > - **运行方式**：
      > >
- `go run cmd/part3/main.go` — 运行第三阶段知识演示
> >   - 练习验证：直接 `go build ./practice/...` 检查编译

---

## 🎯 当前进度

| 阶段       | 状态                                   |
|----------|--------------------------------------|
| 第一阶段     | ✅                                    |
| 第二阶段     | ✅                                    |
| **第三阶段** | **🔄 3.2 gRPC 准备中（3.1 网络编程与补充学习 ✅）** |
| 第四阶段     | 待开始                                  |

> **上次会话：2026-08-17** | 3.1 补充学习全部完成（net.ErrClosed / errgroup / LimitListener）；安装 gRPC 工具链（protoc
> 35.1 + 两个插件 + grpc/protobuf 依赖）

---

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
- [x] context 取消传播 — ctx.Done() / ctx.Err()
- [x] net.ErrClosed + errors.Is — 区分正常关闭与真实错误
- [x] sync.Once — 保证关闭/初始化只执行一次
- [x] golang.org/x/net/netutil.LimitListener — 限制并发连接数

> 说明：以上 6 项已全部完成；net.ErrClosed / errgroup / LimitListener 在 practice/net_extra.go 中实现并验证。

### 🚀 3.2 gRPC 准备清单（2026-08-17 工具链已就绪）

- **目录约定（2026-08-17 拆分演示/练习）**：
    - `proto/demo/hello.proto` — 演示 proto（导师写），生成到 `api/demo/`
    - `proto/calc/calc.proto` — 练习 proto（用户写），生成到 `api/calc/`
    - `api/` — protoc 生成代码（不要手改）
    - `cmd/grpc-demo/main.go` — 快捷运行演示：`go run cmd/grpc-demo/main.go`

- 已安装工具链（用户本地，无需 sudo）：
    - `protoc` 35.1（/home/composer/go/bin/protoc）
    - `protoc-gen-go` 1.36.12
    - `protoc-gen-go-grpc` 1.6.2
- go.mod 依赖：`google.golang.org/grpc v1.83.0`、`google.golang.org/protobuf v1.36.12` 已 tidy 为直接依赖
- 注意：`/home/composer/go/bin` 不在默认 PATH，运行 protoc 前需 `export PATH=$PATH:/home/composer/go/bin` 或使用全路径
- 建议练习（由浅入深）：
    - 3.2-A：定义 `.proto`，生成 Go 代码
    - 3.2-B：Unary RPC（echo / add）
    - 3.2-C：Server-streaming RPC
    - 3.2-D：Client-streaming RPC
    - 3.2-E：Bidirectional-streaming RPC
    - 3.2-F：metadata / 超时 / 拦截器（可选）
- 3.1 补充学习已全部完成，可直接进入 3.2-A

---

## 第三阶段：分布式系统专项 🎯

- [x] 3.1 网络编程（TCP/HTTP Server/Client）
- [ ] 3.2 gRPC（protobuf + Unary/Stream RPC）🔄
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
