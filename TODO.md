# Go 学习路线图

> 基于 [《The Way To Go》](https://learnku.com/docs/the-way-to-go) 教程，目标：为学习**分布式系统**打下坚实基础。

---

> ⚠️ **AI 助手指令（每次会话必读）**
>
> 这是一个 Go 学习项目，用户正在按路线图自学。你的角色是**导师**，不是代劳者。
>
> - **禁止直接实现 TODO/练习**：`practice/exercises.go` 中的 `[TODO]` 标记表示用户尚未完成的练习，不要替用户写代码。
> - **可以做的事**：解释概念、分析已有代码、回答疑问、指出错误、给出思路提示、运行验证命令（如 `go run ./cmd/phase2/`）。
> - **如果用户卡住**：先给思路提示，不要直接给完整答案。用户要求看答案时再展示。
> - **适当扩展学习**：当发现用户缺少某项前置知识（如文件 I/O、JSON 库），主动建议插入补充模块。
> - **检查练习**：按照用户要求检查相关习题。
> - **关键文件**：
    >
- `TODO.md` — 学习路线图和进度
>   - `review/*.go` — 知识点讲解代码
>   - `practice/exercises.go` — 待实现的练习题
>   - `practice/check.go` — 部分练习的自测验证
>   - `cmd/phase2/` — 第二阶段入口

---

## 🎯 当前进度：第二阶段 / 模块一（接口）

> **上次会话结束点: 2026-07-12**

| 阶段                | 状态                       |
|-------------------|--------------------------|
| 第一阶段：复习与查漏补缺      | ✅ 完成                     |
| **第二阶段：接口、测试、并发** | **🔄 文件 I/O & JSON 学习中** |
| 第三阶段：分布式系统专项      | 待开始                      |
| 第四阶段：综合项目         | 待开始                      |

### 补充模块：文件 I/O 与 JSON 基础

> 插入原因：完成 FileStore（练习16）和 JSON tag（练习24）的前置依赖

- [ ] 阅读 `review/io_json.go` — 跑 `go run ./cmd/phase2/`
- [ ] 练习23: JSON 读写（SaveConfig / LoadConfig）
- [ ] 回到练习16: 完成 FileStore（`os.ReadFile` + `json.Marshal` + `os.WriteFile`）
- [ ] 阅读 `review/structs.go` — 结构体嵌入与标签

### 第二阶段的三个模块（按顺序学）

```
接口 (2.1-2.3)  →  测试 (2.4-2.5)  →  并发 (2.6-2.8)
    ↓                    ↓                   ↓
  抽象能力            质量保证            分布式核心
```

### 接口模块学习清单

- [x] 阅读 `review/interfaces.go` — 跑 `go run ./cmd/phase2/`
- [x] 练习15: 实现 sort.Interface 按字符串长度排序
- [x] 练习16: KVStore 接口 + 内存/文件双实现（MemStore ✅, FileStore ⏸️ 等学完 JSON/文件）
- [x] 练习17: 类型断言遍历 Animal 切片
- [x] 练习18: Plugin 模式 — Pipeline 串联处理
- [x] 练习19: 自定义错误类型（实现 error 接口）
- [x] 练习20: 实现 fmt.Stringer 接口（IPAddr）
- [x] 练习21: nil 接口陷阱理解
- [x] 练习22: 接口组合（Printer + Scanner → AllInOne）

---

## 已掌握 ✅

| 章节  | 内容                                      | 状态    |
|-----|-----------------------------------------|-------|
| Ch4 | 基本类型、变量、常量、字符串、指针                       | ✅     |
| Ch5 | if/else、switch、for、break/continue、label | ✅     |
| Ch6 | 函数、闭包、defer、变长参数、多返回值                   | ✅     |
| Ch7 | 数组、切片、make、append、copy                  | ✅     |
| Ch8 | Map、comma-ok 模式、delete                  | ✅     |
| Ch9 | 包管理、go mod                              | ✅     |
| —   | 泛型 (Generics)                           | ✅ 已超前 |
| —   | Functional Options / Builder 模式         | ✅ 已超前 |

---

## 第一阶段：复习与查漏补缺 🔄 ✅ 已完成

> 巩固已有知识，填补遗漏。

- [x] **1.1 复习基础类型与运算符** (`review/basics.go`)
  - [x] 位运算符 (`&`, `|`, `^`, `&^`, `<<`, `>>`)
  - [x] 类型别名 vs 类型定义 (`type A = B` vs `type A B`)
  - [x] `iota` 枚举
- [x] **1.2 复习控制结构** (`review/flow.go`)
  - [x] switch / type switch / label
  - [x] for range 各种写法
- [x] **1.3 复习函数与闭包** (`review/functions.go`)
  - [x] 闭包陷阱、defer LIFO、panic/recover
  - [x] 具名返回值与 defer 交互
- [x] **1.4 复习数组与切片** (`review/collections.go`)
  - [x] 切片底层三元组 (ptr/len/cap)、扩容、共享陷阱
  - [x] 三索引切片、copy 语义
- [x] **1.5 复习 Map** (`review/collections.go`)
  - [x] comma-ok、排序遍历、词频统计
- [x] **1.6 复习错误处理** (`review/functions.go`)
  - [x] errors.Is（哨兵值比较） vs errors.As（类型提取）
  - [x] %w 错误包装链
- [x] **1.7 练习** (`practice/exercises.go`) — 9/9 ✅
  - [x] 练习1: 反转字符串（支持中文）
  - [x] 练习2: 两数之和 O(n) map
  - [x] 练习3: 词频统计 + TopN 排序
  - [x] 练习4: 环形缓冲区 RingBuffer
  - [x] 练习5: 斐波那契记忆化（闭包+map缓存）
  - [x] 练习6: 滑动窗口最大值
  - [x] 练习7: 合并两个有序切片
  - [x] 练习8: 切片去重（保持顺序）
  - [x] 练习9: 展平嵌套切片

---

## 第二阶段：必学核心 🆕

> 这些是 Go 分布式开发的基石。**预计 7-10 天。**

### 接口与反射 (Ch10-Ch11)

- [ ] **2.1 结构体与方法** — [Ch10](https://learnku.com/docs/the-way-to-go/101-structure-definition/3638)
    - [ ] 结构体定义、工厂方法、标签(tag) — tag 在 JSON/protobuf 编解码中无处不在
    - [ ] 匿名字段与结构体嵌入（组合 > 继承）
  - [x] 方法：值接收者 vs 指针接收者 — 接口练习中已掌握
  - [x] `String()` 方法实现自定义格式化 — 练习20已掌握
  - [ ] 练习24: Database 嵌入 Config 结构体
  - [ ] 练习25: APIConfig 带 JSON tag
  - [ ] 练习26: 给 MemStore 添加 String() 方法

- [x] **2.2 接口** — [Ch11](https://learnku.com/docs/the-way-to-go/what-is-the-111-interface/3647) ✅ 已完成
  - [x] 接口即契约：隐式实现（Go 最大特色）
  - [x] 空接口 `interface{}` / `any` 的使用场景
  - [x] 类型断言 `v, ok := x.(T)` 和 type switch
  - [x] 接口嵌套组合
  - [x] 常用标准库接口：`io.Reader`、`io.Writer`、`fmt.Stringer`、`error`、`sort.Interface`

- [ ] **2.3 反射** — [11.10](https://learnku.com/docs/the-way-to-go/1110-reflector/3656)
    - [ ] `reflect.Type` 和 `reflect.Value`
    - [ ] 通过反射读取/修改结构体字段
    - [ ] 理解反射的性能代价与使用场景（ORM、RPC 框架的基础）
    - [ ] 练习：写一个 `StructToMap(s interface{}) map[string]interface{}` 函数

### 错误处理与测试 (Ch13)

- [ ] **2.4 错误处理进阶** — [Ch13.1-13.6](https://learnku.com/docs/the-way-to-go/131-error-handling/3674)
  - [x] `errors.Is`、`errors.As`、`fmt.Errorf("%w", err)` — 1.6 中已复习
  - [x] 自定义错误类型 — 练习19已掌握
    - [ ] `panic` 和 `recover`：何时用、何时不用
    - [ ] 练习：实现一个带超时和重试机制的 HTTP 请求函数

- [ ] **2.5 单元测试
  ** — [Ch13.7-13.9](https://learnku.com/docs/the-way-to-go/unit-testing-and-benchmarking-in-137-go/3680)
    - [ ] `go test`、`TestXxx(t *testing.T)` 基本用法
    - [ ] 表驱动测试 (Table-Driven Tests) — Go 社区标准
    - [ ] 子测试 `t.Run()`
    - [ ] 测试辅助函数 `t.Helper()`
    - [ ] 基准测试 `BenchmarkXxx(b *testing.B)`
    - [ ] 练习：为之前写的 `KVStore` 写完整的表驱动测试 + 基准测试

### 文件与网络 I/O (Ch12)

- [ ] **2.6 文件读写** — [Ch12](https://learnku.com/docs/the-way-to-go/122-file-reading-and-writing/3662)
    - [ ] `os.Open` / `os.Create` / `os.WriteFile` / `os.ReadFile`
    - [ ] `bufio` 缓冲读写（大文件高效处理）
    - [ ] `io.Copy` 流式拷贝
    - [ ] defer 关闭文件的最佳实践
    - [ ] 练习：实现一个日志文件轮转器 (log rotator)

- [ ] **2.7 序列化与编码**
    - [ ] JSON 编解码进阶 (`json.Marshal`/`Unmarshal` 的所有细节)
    - [ ] 自定义 JSON 序列化 (`json.Marshaler`/`json.Unmarshaler`)
    - [ ] Protocol Buffers 入门 — 分布式系统通信标配
    - [ ] 练习：用 JSON 实现一个简单的 Key-Value 持久化存储

### 并发编程 (Ch14) — 🔥 分布式核心

- [ ] **2.8 Goroutine 与 Channel
  ** — [Ch14.1-14.5](https://learnku.com/docs/the-way-to-go/141-concurrency-parallel-and-co-process/3685)
    - [ ] `go` 关键字：goroutine 不是线程，是用户态轻量协程
    - [ ] channel：无缓冲 vs 有缓冲，方向 (`<-chan` / `chan<-`)
    - [ ] `select` 多路复用：超时控制、非阻塞操作、default 分支
    - [ ] channel 关闭：`close()`，range 遍历，`_, ok` 检测
    - [ ] 练习：实现 **生产者-消费者** 模型（多个 producer，多个 consumer）

- [ ] **2.9 并发模式
  ** — [Ch14.8-14.15](https://learnku.com/docs/the-way-to-go/implementation-of-148-inert-generator/3692)
    - [ ] Pipeline 模式：串联多个 stage，每个 stage 用 goroutine + channel
    - [ ] Fan-out / Fan-in 模式：多 worker 并行处理
    - [ ] Futures 模式：提前返回一个 future，异步填充结果
    - [ ] 练习：实现一个 **并发爬虫**（给定 URL 列表，并发抓取，限制并发数）

- [ ] **2.10 同步原语** — [Ch9.3 sync 包](https://learnku.com/docs/the-way-to-go/93-locks-and-sync-packages/3627)
    - [ ] `sync.Mutex` / `sync.RWMutex`：保护共享数据
    - [ ] `sync.WaitGroup`：等待一组 goroutine 完成
    - [ ] `sync.Once`：单次初始化（单例模式，分布式连接池初始化）
    - [ ] `sync/atomic`：原子操作计数器
    - [ ] 练习：实现一个 **线程安全的计数器**（分别用 Mutex 和 atomic 实现，benchmark 对比）

- [ ] **2.11 Context 包** — ⚠️ 教程未涉及，但分布式必学
    - [ ] `context.Background()` / `context.TODO()`
    - [ ] `context.WithCancel` — 取消信号传播
    - [ ] `context.WithDeadline` / `context.WithTimeout` — 超时控制
    - [ ] `context.WithValue` — 传递元数据（trace ID 等）
    - [ ] Context 在 RPC/HTTP 中的传递链路
    - [ ] 练习：实现一个 **模拟微服务调用链**，parent 取消后所有子调用自动取消

---

## 第三阶段：分布式系统专项 🎯

> 学完上面基础后再开始。**预计 10-14 天。**

- [ ] **3.1 网络编程** — [Ch15](https://learnku.com/docs/the-way-to-go/151-tcp-server/3703)
    - [ ] TCP Server/Client：`net.Listen`、`net.Dial`
    - [ ] HTTP Server：`net/http` 标准库，中间件模式
    - [ ] HTTP Client：请求构造、超时、重试、连接池
    - [ ] 练习：实现一个 **HTTP API 网关**（转发请求、超时控制、限流）

- [ ] **3.2 gRPC**
    - [ ] protobuf 定义服务接口
    - [ ] Unary RPC 和 Stream RPC（Server/Client/Bidirectional）
    - [ ] gRPC interceptors（中间件）
    - [ ] 练习：实现一个简单的 **gRPC 键值存储服务**（Put/Get/Delete + Stream 批量操作）

- [ ] **3.3 服务注册与发现**
    - [ ] 理解服务发现原理（client-side / server-side discovery）
    - [ ] 使用 etcd 或 consul 的 Go SDK 实现注册/发现
    - [ ] 健康检查与心跳机制
    - [ ] 练习：实现一个带健康检查的 **简易服务注册中心**

- [ ] **3.4 负载均衡**
    - [ ] 算法：Round Robin、加权轮询、最少连接、一致性哈希
    - [ ] 用 channel 实现一个负载均衡器
    - [ ] 练习：实现一个支持多种策略的 **负载均衡器**，写入 TODO.md 中已有接口的组合

- [ ] **3.5 分布式共识初探**
    - [ ] Raft 算法概念理解（Leader Election, Log Replication, Safety）
    - [ ] 使用 etcd/raft 或 hashicorp/raft 库搭建一个简单的 Raft 集群
    - [ ] 练习：3 节点 Raft 集群，实现一个简单的分布式计数器

- [ ] **3.6 消息队列**
    - [ ] 用 channel 实现内存消息队列
    - [ ] 了解 Kafka / NATS / Redis Stream 基本概念
    - [ ] 练习：实现一个支持 **Pub/Sub 模式** 的内存消息代理（topic 订阅、消息持久化可选）

- [ ] **3.7 分布式追踪与可观测性**
    - [ ] OpenTelemetry Go SDK 集成
    - [ ] Trace 传播（Context + HTTP Header）
    - [ ] 结构化日志 (`slog` 包，Go 1.21+)
    - [ ] 练习：为 gRPC 服务添加 trace 和结构化日志

---

## 第四阶段：综合项目 🏗️

> 将前面所学整合。**预计 7-14 天。**

- [ ] **4.1 分布式 KV 存储**
    - [ ] 单机版 → 支持 gRPC 远程访问
    - [ ] 带 leader/follower 复制
    - [ ] 客户端负载均衡
    - [ ] 支持 Watch 机制（用 channel 实现变更通知）

- [ ] **4.2 简易任务调度系统**
    - [ ] HTTP API 提交任务
    - [ ] Worker pool 执行（goroutine + channel）
    - [ ] 任务状态追踪、重试、超时
    - [ ] 使用 MySQL 存储任务（go-sql-driver/mysql 已在 go.mod 中）

---

## 每日练习清单（贯穿整个学习过程）

| 编号  | 题目                                     | 涉及的技能点                           |
|-----|----------------------------------------|----------------------------------|
| P1  | 反转字符串（支持中文）                            | `[]rune`、Unicode                 |
| P2  | 两数之和（用 map 优化 O(n²)→O(n)）              | Map 查找                           |
| P3  | 实现一个阻塞队列 (BlockingQueue)               | channel、goroutine                |
| P4  | 实现并发安全的 LRU Cache                      | Map + Mutex + 链表                 |
| P5  | 实现超时重试 HTTP Client                     | Context、http.Client、error 处理     |
| P6  | 解析命令行参数写一个简易 CLI 工具                    | `flag` 包、`os.Args`               |
| P7  | 实现 Worker Pool（固定数量 goroutine 处理 jobs） | goroutine、channel、sync.WaitGroup |
| P8  | 用 select 实现超时控制的多数据源聚合查询               | select、time.After                |
| P9  | 实现一个简单的 TCP 代理 (Proxy)                 | `net`、`io.Copy`、goroutine        |
| P10 | 文件内容搜索工具（类似简化版 grep）                   | 文件 I/O、`bufio`、正则                |

---

## 推荐额外资源

- **并发**: [Go Concurrency Patterns (Google I/O 2012)](https://www.youtube.com/watch?v=f6kdp27TYZs)
- **Raft
  **: [Raft 可视化](https://raft.github.io/) + [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
- **分布式**: [MIT 6.824 (现 6.5840) Distributed Systems](https://pdos.csail.mit.edu/6.824/)
- **gRPC**: [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
- **测试**: [Go Testing 最佳实践](https://go.dev/doc/tutorial/add-a-test)

---

## 学习节奏建议

```
Week 1:   第一阶段 1.1-1.5（复习查漏）
Week 2-3: 第二阶段 2.1-2.5（接口、测试、错误处理）
Week 4-5: 第二阶段 2.6-2.11（I/O、并发、Context）
Week 6-8: 第三阶段 3.1-3.7（分布式专项）
Week 9-10: 第四阶段 4.1-4.2（综合项目）
```
