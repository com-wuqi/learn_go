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
> - **当前进度快照（2026-08-30）**：
> >
>   - 第三阶段 3.1 网络编程与补充学习已完成 ✅；3.2 gRPC 全部完成（A-F）；3.2-S gRPC
      补充学习全部完成（status/重试/resolver/TLS/keepalive/连接状态机/流式拦截器/压缩限长/health/反射）
>   - 补充学习已全部完成：errgroup / WaitGroup / context 取消 / net.ErrClosed / sync.Once / LimitListener
>   - 3.3 服务注册与发现（etcd）：已启动（2026-08-30）——3.3-A etcdctl 实操完成 ✅；3.3-B Go client v3
      练习已布置（practice/etcd_basic.go 待完成）；环境：Docker 容器 learn-etcd（v3.5.33）运行中，endpoint=127.0.0.1:2379
> - **项目约定（用户 2026-08-11 明确）**：
> >
>   - 练习题解答写在 `practice/*.go`（用户自己动手）；`review/*.go` 只做大范围复习讲解
>   - 用户偏好流程：演示 → 练习 → 检查 → 修正；练习给函数签名 + TODO 提示，不用填空式 `____`
>   - 检查练习：先读代码 → `go build ./...` → 必要时写临时测试（`practice/zz_check_test.go`，测完删除）
> - **环境/工具备忘**：
> >
>   - 沙箱 bwrap 故障（2026-08-17 实测）：`apply_patch` 新增文件可用，但修改/删除已有文件会报 bwrap 错误；用提权
      `exec_command` + Python 定点替换文件，改完 `gofmt -w` 并编译验证
>   - etcd 环境（2026-08-30 实测）：官方镜像默认只监听容器内 127.0.0.1:2379，且默认 initial-cluster 与探测到的 peer
      地址不一致会启动失败；正确启动命令：
>    
      `docker run -d --name learn-etcd -p 2379:2379 -p 2380:2380 quay.io/coreos/etcd:v3.5.33 /usr/local/bin/etcd --advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://127.0.0.1:2380 --initial-cluster default=http://127.0.0.1:2380`
>   - etcd Go 依赖：`go.etcd.io/etcd/client/v3 v3.5.33`（go get 由用户手动执行成功；proxy.golang.org 不可达，goproxy.cn 可用）
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
>   - `review/grpc_proto.go` — proto 序列化演示
>   - `review/grpc_unary.go` — Unary RPC 演示
>   - `review/grpc_stream.go` — Server-streaming 演示
>   - `review/grpc_client_stream.go` — Client-streaming 演示
>   - `practice/grpc.go` — 3.2-B Unary 练习
>   - `practice/grpc_stream.go` — 3.2-C Server-streaming 练习
>   - `practice/grpc_client_stream.go` — 3.2-D Client-streaming 练习（TODO 待完成）
>   - `review/grpc_bidi.go` / `review/grpc_advanced.go` — 3.2-E/F 演示
>   - `review/grpc_status.go` / `review/grpc_retry.go` / `review/grpc_resolver.go` — 3.2-S 批1 演示
>   - `practice/grpc_bidi.go` / `practice/grpc_advanced.go` — 3.2-E/F 练习
>   - `practice/grpc_extra.go` — 3.2-S 练习（status/重试）
>   - `review/grpc_tls.go` / `review/grpc_keepalive.go` / `review/grpc_connstate.go` — 3.2-S 批2 演示
>   - `practice/grpc_net.go` — 3.2-S 批2 练习（TLS/连接状态机）
>   - `review/grpc_stream_interceptor.go` / `review/grpc_sizes.go` / `review/grpc_health.go` /
      `review/grpc_reflection.go` — 3.2-S 批3 演示
>   - `practice/grpc_tools.go` — 3.2-S 批3 练习（流式拦截器/压缩/health/反射）
>   - `cmd/grpc-demo/` — gRPC 演示快捷运行：`go run cmd/grpc-demo/main.go`
>   - `review/etcd_basic.go` / `cmd/etcd-demo/` — 3.3-B etcd 基础演示：`go run cmd/etcd-demo/main.go`
>   - `practice/etcd_basic.go` — 3.3-B 练习（Go client v3 基础，TODO 待完成）
>   - `cmd/part1/` / `cmd/part2/` / `cmd/part3/` — 按阶段运行演示
> > - **运行方式**：
>   >
- `go run cmd/part3/main.go` — 运行第三阶段知识演示
> >   - 练习验证：直接 `go build ./practice/...` 检查编译

---

## 🎯 当前进度

| 阶段         | 状态                                                        |
|--------------|-------------------------------------------------------------|
| 第一阶段     | ✅                                                          |
| 第二阶段     | ✅                                                          |
| **第三阶段** | **🔄 3.3 服务注册与发现（etcd）进行中（3.3-A ✅，3.3-B 练习待完成）** |
| 第四阶段     | 待开始                                                      |

> **上次会话：2026-08-30** | 3.3 启动：环境（learn-etcd 容器 v3.5.33）+ 3.3-A etcdctl 演示 + 3.3-B Go
基础演示（review/etcd_basic.go）完成；3.3-B 练习已布置（practice/etcd_basic.go）

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
- [x] context 取消传播 — ctx.Done () / ctx.Err ()
- [x] net.ErrClosed + errors.Is — 区分正常关闭与真实错误
- [x] sync.Once — 保证关闭/初始化只执行一次
- [x] golang.org/x/net/netutil.LimitListener — 限制并发连接数

> 说明：以上 6 项已全部完成；net.ErrClosed / errgroup / LimitListener 在 practice/net_extra.go 中实现并验证。

### 🚀 3.2 gRPC 进度与关键结论

- **目录约定（2026-08-17 拆分演示/练习）**：
    - `proto/demo/hello.proto` — 演示 proto（导师写），生成到 `api/demo/`
    - `proto/calc/calc.proto` — 练习 proto（用户写），生成到 `api/calc/`
    - `api/` — protoc 生成代码（不要手改）
    - `cmd/grpc-demo/main.go` — 快捷运行演示：`go run cmd/grpc-demo/main.go`
  - 重新生成：
    `protoc -I proto --go_out=. --go_opt=module=LearnGo --go-grpc_out=. --go-grpc_opt=module=LearnGo proto/xxx.proto`（先
    export PATH）

- 已安装工具链（用户本地，无需 sudo）：
    - `protoc` 35.1（/home/composer/go/bin/protoc）
    - `protoc-gen-go` 1.36.12
    - `protoc-gen-go-grpc` 1.6.2
- go.mod 依赖：`google.golang.org/grpc v1.83.0`、`google.golang.org/protobuf v1.36.12` 已 tidy 为直接依赖
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
  - 连接状态机：无在途请求断连回 IDLE，有在途请求才进 TRANSIENT_FAILURE；默认 Unary 是 fail-fast，WaitForReady (true) 才会等重连
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

### 🚀 3.3 服务注册与发现（etcd）

> 目标：掌握 etcd 核心机制（KV / 前缀 / Watch / Lease），并实现服务注册 + 发现。

- [x] 3.3-A：etcd 概念 + etcdctl 实操（put/get/--prefix/lease TTL/expire）
- [ ] 3.3-B：Go client v3 基础（Put/Get/ListPrefix/Watch/LeasePut）— practice/etcd_basic.go
- [ ] 3.3-C：服务注册（lease 注册 + KeepAlive 自动续约 + 优雅注销）
- [ ] 3.3-D：服务发现（Watch 维护本地实例列表 + 故障容错）
- [ ] 3.3-E：综合——多副本 gRPC 服务注册到 etcd，客户端经 etcd 发现并调用

> 3.3-B 关键结论（2026-08-30 沉淀）：
> - key 组织：`/services/<服务名>/<实例ID>` -> 地址，前缀查询即"发现该服务的所有实例"
> - Watch 是增量事件流（PUT/DELETE），本地缓存 + watch = 服务发现的实时状态
> - Lease 是"心跳"：KeepAlive 续租保持存活，租约过期/Revoke 时绑定 key 自动删除（实例下线自愈）
> - clientv3.New 需显式 DialTimeout；连接本身是 lazy 的，第一个 RPC 才真正拨号

## 第三阶段：分布式系统专项 🎯

- [x] 3.1 网络编程（TCP/HTTP Server/Client）
- [x] 3.2 gRPC（protobuf + Unary/Stream RPC）
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
