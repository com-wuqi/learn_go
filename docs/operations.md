# 操作与文件索引

> 本文件由 README.md 拆分而来（2026-09-07）。README 只保留索引与进度；需要环境备忘、文件导航、运行方式时查阅本文档。

## 环境/工具备忘

- 沙箱 bwrap 故障（2026-08-17 实测）：`apply_patch` 新增文件可用，但修改/删除已有文件会报 bwrap 错误；用提权
  `exec_command` + Python 定点替换文件，改完 `gofmt -w` 并编译验证
- etcd 环境（2026-08-30 实测）：官方镜像默认只监听容器内 127.0.0.1:2379，且默认 initial-cluster 与探测到的 peer
  地址不一致会启动失败；正确启动命令：

  `docker run -d --name learn-etcd -p 2379:2379 -p 2380:2380 quay.io/coreos/etcd:v3.5.33 /usr/local/bin/etcd --advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://127.0.0.1:2380 --initial-cluster default=http://127.0.0.1:2380`
- etcd Go 依赖：`go.etcd.io/etcd/client/v3 v3.5.33`（go get 由用户手动执行成功；proxy.golang.org 不可达，goproxy.cn 可用）
- gRPC 工具链（用户本地 /home/composer/go/bin，不在默认 PATH）：protoc 35.1、protoc-gen-go 1.36.12、protoc-gen-go-grpc
  1.6.2；使用前 export PATH 或全路径
- 环境有 `HTTP_PROXY=http://127.0.0.1:7897/` 且 `NO_PROXY` 不含 `[::]`：本地 HTTP 演示客户端地址用 `127.0.0.1`（绑
  `:0` 会被代理返回 502）
- `review/networking.go` FlusherDemo 有一行 `} // codex resume <uuid>` 残留注释（未确认是否删除）；UploadDemo 的
  `ParseMultipartForm(4 << 20)` 是用户有意改的（演示大文件落盘）

## 关键文件索引

- `review/basics.go` … `review/http.go` — 一、二阶段知识点讲解
- `review/networking.go` — 三阶段网络编程知识点讲解
- `practice/exercises.go` — 练习 1-9（基础）
- `practice/interfaces.go` — 练习 15-22（接口）
- `practice/io.go` — 练习 23（文件/JSON）
- `practice/structs.go` — 练习 24-30（结构体）
- `practice/reflect.go` — 练习 31-33, 36（反射 + panic）
- `practice/kvstore_test.go` — 练习 34-35（测试）
- `practice/time.go` — 练习 37-38（time 包）
- `practice/io_reader.go` — 练习 39-41（io.Reader/Writer）
- `practice/http.go` — 练习 42-44（HTTP）
- `practice/goroutine.go` — 练习 45-46（goroutine/channel）
- `practice/concurrency.go` — 练习 10-14, 47-53（并发综合）
- `practice/review.go` — 复习题 A-M（综合）
- `practice/networking.go` — 练习 N-V + U-1（网络编程）
- `practice/net_extra.go` — 3.1 补充练习（net.ErrClosed / errgroup / LimitListener）
- `review/net_extra.go` — 3.1 补充知识讲解（net.ErrClosed / errgroup / LimitListener）
- `review/grpc_proto.go` — proto 序列化演示
- `review/grpc_unary.go` — Unary RPC 演示
- `review/grpc_stream.go` — Server-streaming 演示
- `review/grpc_client_stream.go` — Client-streaming 演示
- `practice/grpc.go` — 3.2-B Unary 练习
- `practice/grpc_stream.go` — 3.2-C Server-streaming 练习
- `practice/grpc_client_stream.go` — 3.2-D Client-streaming 练习（已完成）
- `review/grpc_bidi.go` / `review/grpc_advanced.go` — 3.2-E/F 演示
- `review/grpc_status.go` / `review/grpc_retry.go` / `review/grpc_resolver.go` — 3.2-S 批1 演示
- `practice/grpc_bidi.go` / `practice/grpc_advanced.go` — 3.2-E/F 练习
- `practice/grpc_extra.go` — 3.2-S 练习（status/重试）
- `review/grpc_tls.go` / `review/grpc_keepalive.go` / `review/grpc_connstate.go` — 3.2-S 批2 演示
- `practice/grpc_net.go` — 3.2-S 批2 练习（TLS/连接状态机）
- `review/grpc_stream_interceptor.go` / `review/grpc_sizes.go` / `review/grpc_health.go` /
  `review/grpc_reflection.go` — 3.2-S 批3 演示
- `practice/grpc_tools.go` — 3.2-S 批3 练习（流式拦截器/压缩/health/反射）
- `cmd/grpc-demo/` — gRPC 演示快捷运行：`go run cmd/grpc-demo/main.go`
- `review/etcd_basic.go` / `cmd/etcd-demo/` — 3.3-B etcd 基础演示：`go run cmd/etcd-demo/main.go`
- `practice/etcd_basic.go` — 3.3-B 练习（Go client v3
  基础，已完成：ConnectEtcd/EtcdPutGet/EtcdListPrefix/EtcdWatchPrefix/EtcdLeasePut）
- `cmd/part1/` / `cmd/part2/` / `cmd/part3/` — 按阶段运行演示

## 运行方式

- `go run cmd/part3/main.go` — 运行第三阶段知识演示
    - 练习验证：直接 `go build ./practice/...` 检查编译
