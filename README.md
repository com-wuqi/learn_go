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
> - **检查练习**：按用户要求检查相关习题。流程：读代码 → `go build ./...` → 必要时写临时测试（`practice/zz_check_test.go`
    ，测完删除）。
> - **项目约定（用户 2026-08-11 明确）**：
>   - 练习题解答写在 `practice/*.go`（用户自己动手）；`review/*.go` 只做大范围复习讲解
>   - 用户偏好流程：演示 → 练习 → 检查 → 修正；练习给函数签名 + TODO 提示，不用填空式 `____`
> - **工作区文档结构（2026-09-07 拆分，避免 README 过长）**：
>   - `docs/operations.md` — 环境/工具备忘、关键文件索引、运行方式（涉及工具/文件导航时打开）
>   - `docs/notes-3.1-networking.md` / `docs/notes-3.2-grpc.md` / `docs/notes-3.3-etcd.md` — 各模块学习笔记（对应模块进行中/检查时打开）
> - **当前进度快照（2026-09-07）**：3.1 ✅、3.2（含 3.2-S）✅、3.3 进行中（3.3-A/B ✅，3.3-C 待开始）

---

## 🎯 当前进度

| 阶段         | 状态                                                                        |
|--------------|-----------------------------------------------------------------------------|
| 第一阶段     | ✅                                                                          |
| 第二阶段     | ✅                                                                          |
| **第三阶段** | **🔄 3.3 服务注册与发现（etcd）进行中（3.3-A ✅，3.3-B ✅，3.3-C 待开始）** |
| 第四阶段     | 待开始                                                                      |

> **上次会话：2026-09-07** | 3.3-B 练习 1-5 检查通过（行为测试 + watch goroutine 防泄漏验证）；README 拆分完成（索引 +
> docs/operations.md + 按模块笔记）
>
> **下个会话起点**：3.3-C 服务注册——先讲 key 设计与 KeepAlive 生命周期/优雅注销，再建 practice/etcd_register.go。
> 若 learn-etcd 容器已停止，按 docs/operations.md 中的启动命令拉起（endpoint=127.0.0.1:2379）。

---

## 🗂️ 文档导航

- `docs/operations.md` — 环境/工具备忘、关键文件索引、运行方式
- `docs/notes-3.1-networking.md` — 3.1 网络编程：练习 N-V、关键结论、补充学习
- `docs/notes-3.2-grpc.md` — 3.2 gRPC：练习进度、工具链、关键结论、3.2-S 补充学习
- `docs/notes-3.3-etcd.md` — 3.3 服务注册与发现（etcd）：A-E 计划与关键结论

---

## 路线图

### 第三阶段：分布式系统专项 🎯

- [x] [3.1 网络编程（TCP/HTTP Server/Client）](docs/notes-3.1-networking.md)
- [x] [3.2 gRPC（protobuf + Unary/Stream RPC）](docs/notes-3.2-grpc.md)
- [ ] [3.3 服务注册与发现（etcd）](docs/notes-3.3-etcd.md)
- [ ] 3.4 负载均衡
- [ ] 3.5 分布式共识（Raft）
- [ ] 3.6 消息队列（Pub/Sub）
- [ ] 3.7 可观测性（OpenTelemetry）

### 第四阶段：综合项目 🏗️

- [ ] 4.1 分布式 KV 存储
- [ ] 4.2 任务调度系统

## 推荐资源

- **并发**: [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- **Raft**: [Raft 可视化](https://raft.github.io/)
- **分布式**: [MIT 6.824](https://pdos.csail.mit.edu/6.824/)
- **gRPC**: [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
