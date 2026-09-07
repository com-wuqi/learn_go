# 3.3 服务注册与发现（etcd）学习笔记

> 本文件由 README.md 拆分而来（2026-09-07）。

### 🚀 3.3 服务注册与发现（etcd）

> 目标：掌握 etcd 核心机制（KV / 前缀 / Watch / Lease），并实现服务注册 + 发现。

- [x] 3.3-A：etcd 概念 + etcdctl 实操（put/get/--prefix/lease TTL/expire）
- [x] 3.3-B：Go client v3 基础（Put/Get/ListPrefix/Watch/LeasePut）— practice/etcd_basic.go
- [ ] 3.3-C：服务注册（lease 注册 + KeepAlive 自动续约 + 优雅注销）
- [ ] 3.3-D：服务发现（Watch 维护本地实例列表 + 故障容错）
- [ ] 3.3-E：综合——多副本 gRPC 服务注册到 etcd，客户端经 etcd 发现并调用

> 3.3-B 关键结论（2026-08-30 沉淀）：
> - key 组织：`/services/<服务名>/<实例ID>` -> 地址，前缀查询即"发现该服务的所有实例"
> - Watch 是增量事件流（PUT/DELETE），本地缓存 + watch = 服务发现的实时状态
> - Lease 是"心跳"：KeepAlive 续租保持存活，租约过期/Revoke 时绑定 key 自动删除（实例下线自愈）
> - clientv3.New 需显式 DialTimeout；连接本身是 lazy 的，第一个 RPC 才真正拨号

> 下一步（3.3-C）：服务注册设计——实例 key 生成、lease + KeepAlive（ctx 生命周期）、优雅注销（Revoke/Close）。
