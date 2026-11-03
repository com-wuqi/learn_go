# 3.3 服务注册与发现（etcd）学习笔记

> 本文件由 README.md 拆分而来（2026-09-07）。

### 🚀 3.3 服务注册与发现（etcd）

> 目标：掌握 etcd 核心机制（KV / 前缀 / Watch / Lease），并实现服务注册 + 发现。

- [x] 3.3-A：etcd 概念 + etcdctl 实操（put/get/--prefix/lease TTL/expire）
- [x] 3.3-B：Go client v3 基础（Put/Get/ListPrefix/Watch/LeasePut）— practice/etcd_basic.go
- [x] 3.3-C：服务注册（lease 注册 + KeepAlive 自动续约 + 优雅注销）— practice/etcd_register.go（2026-10-16 检查通过）
- [x] 3.3-D：服务发现（Watch 维护本地实例列表 + 故障容错）— practice/etcd_discover.go（2026-11-03 验收通过）
- [ ] 3.3-E：综合——多副本 gRPC 服务注册到 etcd，客户端经 etcd 发现并调用
  - [ ] 3.3-E1：ServiceDiscovery 订阅契约（初始快照 + 变更通知 + 取消/关闭）
  - [ ] 3.3-E2：etcd-backed gRPC resolver（初始地址 + 动态更新）
  - [ ] 3.3-E3：多副本注册、round_robin 调用与下线故障切换

> 3.3-B 关键结论（2026-08-30 沉淀）：
> - key 组织：`/services/<服务名>/<实例ID>` -> 地址，前缀查询即"发现该服务的所有实例"
> - Watch 是增量事件流（PUT/DELETE），本地缓存 + watch = 服务发现的实时状态
> - Lease 是"心跳"：KeepAlive 续租保持存活，租约过期/Revoke 时绑定 key 自动删除（实例下线自愈）
> - clientv3.New 需显式 DialTimeout；连接本身是 lazy 的，第一个 RPC 才真正拨号

> 依赖说明（2026-09-20）：go.mod 里的 `github.com/golang/protobuf v1.5.4 // indirect` 无法移除——本项目代码零引用旧包
> （`api/` 生成代码、`practice/`、`review/` 全部用 `google.golang.org/protobuf`），但 etcd client v3 自身仍 import 它
> （`client/v3/kv.go`、`api/v3/etcdserverpb` 等），换到最新 v3.8.0-alpha.0 也一样（上游标注 TODO: remove for a supported
> version）。
> 该模块现在只是转发到 `google.golang.org/protobuf` 的兼容 shim，留在 go.mod 属正常现象，不是本项目漏改。

> 3.3-C 关键结论（2026-10-16 沉淀）：
> - key 用 `/services/<服务名>/<实例ID>`，必须走统一的 ServiceKey，前缀写错发现端就查不到（单数/复数踩过一次）
> - 注册 = Grant + Put (WithLease)；心跳 = KeepAlive (ctx)，ctx 必须从调用方派生且为进程级生命周期
> - KeepServiceAlive 的契约：立即返回；正常取消静默关闭 channel；ctx 未取消而续约断开才算"租约丢失"，要上报错误
> - stop 幂等用 sync.Once 只包住"cancel + Revoke"，不要包住续约循环（Once 是同步执行，会阻塞整个函数）
> - errors.Join (nil, nil) 返回 nil；不要用 fmt.Errorf ("%w", nil) 拼错误（会产出非 nil 的 %!w (<nil>)）
> - 清理动作要用独立 ctx：被 cancel 的父 ctx 派生出的子 ctx 一创建就是已取消的，Revoke 会失败并泄漏租约

> 3.3-D 关键结论（2026-11-03 沉淀）：
> - 一致性起点：先做 linearizable Get 得到全量快照和 `Header.Revision`，再从 `revision+1` 开始 Watch
> - 本地状态：`instanceID -> addr`，后台 goroutine 单写，读取方在锁内复制 map
> - 故障恢复：普通断连交给 etcd client；Watch 被压缩时重新 Get 快照并原子替换本地状态
> - 生命周期：发现器不拥有 client；ctx 取消或 Close 时停止 Watch，Close 必须幂等
> - 两个 revision 必须分开：`appliedRevision` 表示本地状态已应用到哪一版，`nextWatchRevision` 表示下一次从哪一版开始监听
> - errors 采用 buffered + 非阻塞上报；正常 context 取消/Close 不报告为故障，channel 在后台退出时关闭
> - Close 必须 cancel 后等待 done；只 cancel 不等待不满足“停止并清理”的契约
> - 配套演示：`review/etcd_discover.go`（完整状态流 + 真实 `ErrCompacted` 恢复）
> - 稳定验收：`practice/etcd_discover_test.go`；etcd 集成测试通过 `ETCD_ENDPOINT` 启用

> 3.3-E 契约草案（按“契约/测试先行”推进）：
> - 服务侧：两个 Greeter 副本分别监听随机端口，以 lease 注册 `/services/<name>/<instanceID>`；心跳、注销、gRPC Stop 的归属明确
> - 发现侧：复用 ServiceDiscovery 的全量快照 + 增量 Watch，不为 resolver 另写一套轮询发现
> - resolver：初始地址立即注入；发现变化时向 ClientConn 推送新地址；关闭 resolver 后停止订阅且不关闭共享 etcd client
> - 客户端：round_robin 调用两个副本；停止一个副本并注销后，已知地址逐步移除，剩余副本继续服务
> - 测试：先写 resolver 的无网络单元测试和 etcd/gRPC 集成测试，再建立练习骨架；实现细节不写进题目

> 3.3-E1 已建立契约与验收（2026-11-03，待实现）：
> - `Subscribe` 返回独立实例快照、revision 和只读更新信号
> - 订阅使用 buffered=1 + 合并通知；消费者收到信号后重新读取 `Instances`
> - 取消订阅幂等且不影响发现器；发现器退出时关闭全部订阅
> - 先注册订阅再读快照，允许额外一次无害通知，但禁止漏事件
> - 测试：`practice/etcd_subscribe_test.go`
