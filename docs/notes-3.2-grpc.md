# 3.2 gRPC 学习笔记

> 本文件由 README.md 拆分而来（2026-09-07），含 3.2-S 补充学习。

### 🚀 3.2 gRPC 进度与关键结论

- **目录约定（2026-08-17 拆分演示/练习）**：
    - `proto/demo/hello.proto` — 演示 proto（导师写），生成到 `api/demo/`
    - `proto/calc/calc.proto` — 练习 proto（用户写），生成到 `api/calc/`
    - `api/` — protoc 生成代码（不要手改）
    - `cmd/grpc-demo/main.go` — 快捷运行演示：`go run cmd/grpc-demo/main.go`
    - 重新生成：
      `protoc -I proto --go_out=. --go_opt=module=LearnGo --go-grpc_out=. --go-grpc_opt=module=LearnGo proto/xxx.proto`
      （先
      export PATH）

- 已安装工具链（用户本地，无需 sudo）：
    - `protoc` 35.1（/home/composer/go/bin/protoc）
    - `protoc-gen-go` 1.36.12
    - `protoc-gen-go-grpc` 1.6.2
- go.mod 依赖：`google.golang.org/grpc v1.84.0`、`google.golang.org/protobuf v1.36.12` 已 tidy 为直接依赖
- 注意：`/home/composer/go/bin` 不在默认 PATH，运行 protoc 前需 `export PATH=$PATH:/home/composer/go/bin` 或使用全路径
- 建议练习（由浅入深）：
    - [x] 3.2-A：定义 `.proto`，生成 Go 代码
    - [x] 3.2-B：Unary RPC（echo / add）
    - [x] 3.2-C：Server-streaming RPC
    - [x] 3.2-D：Client-streaming RPC
    - [x] 3.2-E：Bidirectional-streaming RPC
    - [x] 3.2-F：metadata / 超时 / 拦截器（可选）
- 3.2-A ~ 3.2-F 全部完成（2026-08-21 check 通过）
- ⚠️ 已解决：calc.proto 已加 Sum/EchoSum RPC，`go build ./...` 通过
- **3.2 关键结论（2026-08-18 沉淀）**：
    - `grpc.NewClient`（及旧 `grpc.Dial`）必须显式传 transport credentials，否则返回 `no transport security set`；本地明文用
      `grpc.WithTransportCredentials(insecure.NewCredentials())`
    - Unary server 方法带 `ctx context.Context`；流式方法不带 ctx，从 `stream.Context()` 取
        - 实现 server 必须嵌入 `UnimplementedXxxServer`（向前兼容）
        - Server-streaming：服务端 `stream.Send`，客户端 `Recv` 直到 `io.EOF`
        - Client-streaming：客户端 `Send` + `CloseAndRecv`，服务端 `Recv` 直到 `io.EOF` 后 `SendAndClose`
        - `GracefulStop`/`Stop` 会关闭 listener，不要再手动 `ln.Close()`（避免 double-close）
        - 流式 Recv 循环里非 `io.EOF` 错误要 `return`，避免访问 nil 响应
- **3.2 补充结论（2026-08-21 沉淀）**：
    - EOF 是正常收尾信号：Recv 循环里 EOF 分支要 `return nil`，当错误返回会让对端看到 Unknown/EOF
    - errgroup 的 join 语义：任何退出路径都要 `Wait()`，否则 goroutine 可能比函数晚一步退出
    - metadata 是 context.WithValue 之上的一层协议化封装：key 自动小写、值限 string、`-bin` 后缀走 base64；WithValue 只做进程内分发
    - 拦截器不注册就是死代码：ChainUnaryInterceptor 挂到 server/client 上才生效；链式用 Chain
    - 双向流 ctx 覆盖整条流而不是每条消息；单条超时需自己 race，且超时后流作废（不能并发 Recv）
    - status：code 分类 + status.FromError + WithDetails 附加任意 proto 消息作为错误详情
    - retryPolicy：只重试 retryableStatusCodes 里的错误；自定义 resolver（Build/UpdateState）+ round_robin 实现多后端分发
- **3.2-S 补充结论（2026-08-30 沉淀，批2/批3）**：
    - 客户端 keepalive 的 Time 有 10s 硬性下限（WithKeepaliveParams 会把更小值静默抬到 10s）；服务端
      EnforcementPolicy.MinTime 违规时发 GOAWAY (ENHANCE_YOUR_CALM, too_many_pings)
    - 连接状态机：无在途请求断连回 IDLE，有在途请求才进 TRANSIENT_FAILURE；默认 Unary 是 fail-fast，WaitForReady (true)
      才会等重连
    - TLS：RootCAs 管信任、Certificates 管 mTLS 身份、ServerName 管校验名；useMTLS 门控，别把零值证书塞进 Certificates
    - 压缩：import gzip 只是注册能力，UseCompressor 才是开关；服务端响应压缩回声请求编码；大小限制按解压后校验；小消息压缩反而变大
    - 流式拦截器每条流（每条 RPC）触发一次，不是每条消息；grpc_health_v1 是标准健康检查协议；反射协议：请求类型决定响应里哪个
      oneof 字段有值

---

### 🧩 3.2-S gRPC 补充学习（2026-08-21 新增）

> 目标：补上 gRPC 生产级缺口 + 分布式前置知识，分三批推进。

- [x] status 错误处理（codes / status.FromError / WithDetails）
- [x] 重试策略（retryPolicy / retryableStatusCodes）
- [x] resolver + 负载均衡（自定义 resolver / round_robin）
- [x] TLS / mTLS（credentials.NewTLS）
- [x] keepalive（长连接保活 / 半开连接探测）
- [x] 连接状态机（ConnectivityState）
- [x] 流式拦截器（ChainStreamInterceptor）
- [x] 压缩 + 消息大小限制
- [x] health 检查（grpc_health_v1）
- [x] 反射 + grpcurl 调试
