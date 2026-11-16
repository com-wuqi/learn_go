# 3.3-R 服务注册与发现复习

> 目的不是再做一组 API 填空，而是把 3.3-B 到 E2 的知识串成一条可解释、可运行、可检查的生命周期。
> E3A 暂缓；完成本文的概念题、代码路线和实验后再继续。

## 1. 全链路心智模型

```text
gRPC 副本
  └─ RegisterServiceWithHeartbeat
       ├─ Grant lease
       ├─ Put /services/<service>/<instanceID> = host:port
       └─ KeepAlive(ctx)

etcd
  ├─ Get prefix @ revision R
  └─ Watch prefix from R+1

ServiceDiscovery
  ├─ 本地 instanceID -> addr
  ├─ 合并变化并通知 subscriber
  └─ ErrCompacted 时全量重同步

Resolver
  ├─ Build 同步发布初始完整地址
  ├─ Updates → UpdateState(完整新列表)
  ├─ Discovery Errors → ClientConn.ReportError
  └─ Close → Cancel subscription + Close discovery

gRPC ClientConn / round_robin
  └─ 把 RPC 分发到当前 READY 的副本
```

先记住“谁创建、谁关闭、谁只使用”：

| 资源                  | 创建者           | 拥有者与释放方式                        |
|-----------------------|------------------|-----------------------------------------|
| etcd client           | 应用入口         | 应用关闭；resolver/discovery 不关闭     |
| gRPC server/listener  | 副本启动代码     | 副本 Stop 关闭                          |
| etcd lease/key        | 注册器           | Stop 时 Revoke；进程崩溃时等待租约过期  |
| ServiceDiscovery      | resolver/应用    | 创建它的组件 Close                      |
| DiscoverySubscription | discovery 调用方 | 调用方 Cancel；discovery 退出时兜底关闭 |
| gRPC ClientConn       | 应用入口         | 应用 Close；其 resolver 随后关闭        |

## 2. 概念题：先不看答案

// 参考答案见下面

1. 为什么服务前缀必须是 `/services/echo/`，而不是 `/services/echo`？

> 确保命名空间正确，防止意外匹配 echo*/

2. 进程正常下线与进程崩溃时，etcd 实例 key 分别怎样消失？

> once revoke; ttl到期

3. `Get` 在 revision `R` 返回快照后，为什么 Watch 要从 `R+1` 开始？

> 前面的变更已经在Get结果里面，watch接下来的变更

4. Watch 收到 `ErrCompacted` 后为什么不能继续补事件？

> 历史版本已经被压缩了，revision无法找回

5. `appliedRevision` 和 `nextWatchRevision` 分别表示什么？

> 已经应用的revision和接下来watch的revision

6. 为什么 subscriber 的 Updates channel 使用容量 1 和非阻塞发送？

> 防止慢消费者阻塞Resolver

7. 为什么 discovery 的错误不能阻塞在 `errors <- err`？

> 会阻塞watch和event

8. 订阅通知到达后，为什么 resolver 要重新读取整张实例表
9. `ResolveNow` 和 `Updates` 都触发发布，它们的语义差别是什么？
10. 为什么 resolver.Close 不能关闭共享 etcd client？
11. 为什么 resolver.Close 不能关闭 refresh channel？
12. 副本 Stop 为什么应先注销 lease，再停止 gRPC server？
13. 哪些后台错误属于“恢复后继续”，哪些表示“生命周期已结束”？
14. 如果一个 etcd transaction 同时修改两个目标实例，revision 和事件应怎样推进？

## 3. 代码阅读路线

按以下顺序阅读，每处只回答“输入、输出、所有权、失败行为”：

1. `practice/etcd_register.go`
    - `RegisterServiceWithHeartbeat`
    - `KeepServiceAlive`
    - stop 的幂等与独立 cleanup context
2. `practice/etcd_discover.go`
    - `LoadServiceSnapshot`
    - `Subscribe`
    - `run`
    - `reportError`
    - `Close`
3. `practice/etcd_resolver.go`
    - `ResolverServiceName`
    - `ResolverAddresses`
    - `Build`
    - `Run`
    - `Close`
4. `review/grpc_resolver.go`
    - `resolver.Builder`、`ClientConn.UpdateState`、`round_robin`
5. `practice/etcd_grpc_replica.go`
    - 先只阅读契约，不实现

## 4. 可运行实验

### 实验 A：观察发现状态流

```bash
go run cmd/etcd-demo/main.go
```

记录：

- 初始快照有哪些实例；
- PUT 新增/覆盖、DELETE 移除时输出如何变化；
- 压缩后的旧 Watch 为什么被取消；
- 重同步后的 revision 起点。

### 实验 B：验证并发和生命周期

```bash
ETCD_ENDPOINT=127.0.0.1:2379 go test -race ./practice -run 'Etcd(Discovery|Resolver)' -v
```

记录：

- 哪些测试验证业务状态；
- 哪些测试验证关闭/并发/所有权；
- 哪个测试如果失败，最可能意味着 goroutine 泄漏。

### 实验 C：解释差异

不看代码，分别写出下面两条路径：

```text
正常网络短暂断开：
进程重启且缓存丢失：
```

再对照 `ServiceDiscovery.Run` 验证你的答案。

## 5. 口述验收

完成实验后，用不超过两分钟回答：

> 一个 gRPC 副本启动后，地址如何进入 etcd；客户端如何得到初始地址；之后新增、修改、删除实例如何更新连接；lease 过期和 Watch
> 压缩分别由谁恢复；最后所有资源按什么顺序关闭？

然后补充说明：

- 为什么不重复应用 revision `R`；
- 为什么错误通道不能阻塞；
- 为什么 Close 必须等待后台结束；
- 为什么 resolver 不拥有共享 etcd client。

## 6. 参考答案

1. 尾随 `/` 把范围限制在同一个服务路径下；否则 `echo` 可能匹配 `echo2` 等兄弟前缀。当前代码还会拒绝包含额外 `/` 的实例 ID。
2. 正常下线调用 Revoke，绑定 key 立即删除；崩溃时 KeepAlive 停止，租约到期后 etcd 自动删除 key。
3. 快照已经包含 revision `R` 及之前的全部状态；从 `R+1` 开始只接收后续事务，避免重复应用 `R`。
4. 压缩后的历史事件已不存在，无法补齐；必须重新 Get 全量快照并原子替换本地状态。
5. 前者是当前本地状态对应版本；后者是下一次 Watch 应请求的首个 revision，通常为 `appliedRevision+1`。
6. 通知只是边缘信号，不携带状态；容量 1 可合并积压，非阻塞发送确保慢消费者不会卡住 discovery。
7. discovery 必须继续维护状态和恢复 Watch；错误消费速度不能反过来阻塞核心状态机。
8. Updates 是“状态可能变化”的信号，不是完整状态；必须在最新本地快照上重新生成完整地址列表。
9. Updates 表示 etcd 状态变化；ResolveNow 是 gRPC 的重新解析提示，etcd 状态未必变化，二者可调用同一个发布函数。
10. client 由应用共享，resolver 只是使用者；关闭它会破坏其他注册器、发现器或客户端。
11. refresh 可能被并发 ResolveNow 发送；关闭后发送会 panic。channel 可不显式关闭，交给垃圾回收。
12. 先注销确保发现端尽快停止把新请求路由到该副本；随后 Stop 只负责终止正在执行的 RPC 和监听器。
13. Watch 压缩可由全量重同步恢复；普通网络断开由 etcd client 重连；Updates/Errors 关闭、context 取消或 Close 表示生命周期结束。
14. 同一 transaction 产生的多个事件共享一个 revision；先应用整个响应的事件批次，再推进 appliedRevision 和
    nextWatchRevision。

## 7. 复习通过后的衔接

回到 `practice/etcd_grpc_replica.go`，先实现 `labeledGreeter` 和 `StartRegisteredGreeter`，再处理 `Stop`。
不要同时开始客户端 round_robin 集成；E3B 会单独验证地址更新和副本下线切换。
