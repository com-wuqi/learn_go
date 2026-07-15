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
    >
- `review/*.go` — 知识点讲解代码（运行入口：`go run ./cmd/phase2/`）
>   - `practice/exercises.go` — 练习 1-9 已掌握
>   - `practice/interfaces.go` — 练习 15-22 已完成
>   - `practice/io.go` — 练习 23 已完成
>   - `practice/structs.go` — 练习 24-30 已完成
>   - `practice/concurrency.go` — 练习 10-14 待学

---

## 🎯 当前进度

| 阶段                | 状态           |
|-------------------|--------------|
| 第一阶段：复习与查漏补缺      | ✅            |
| **第二阶段：接口、测试、并发** | **🔄 测试学习中** |
| 第三阶段：分布式系统专项      | 待开始          |
| 第四阶段：综合项目         | 待开始          |

### 🔥 当前：反射 (Reflection) ✅ 已完成

- [x] 阅读 `review/reflect.go`
- [x] 练习31: StructToMap
- [x] 练习32: FillDefaults
- [x] 练习33: PrintTags

### 当前：测试 & panic/recover

- [ ] 阅读 `review/testing.go` — panic/recover + 测试基础
- [ ] 练习34: 为 KVStore 写表驱动测试
- [ ] 练习35: 为 RingBuffer 写基准测试
- [ ] 练习36: panic/recover — SafeCall

### 补充模块：time 包

> 前置于并发（超时控制、定时任务）和 HTTP 重试都需要

- [ ] `time.Now`、`time.Since`、`time.Sleep`
- [ ] `time.After`、`time.Ticker`、`time.AfterFunc`
- [ ] `time.Duration` 单位转换
- [ ] 练习：带超时的操作、定时轮询

### 补充模块：io.Reader/Writer 实战

> 已掌握接口定义，缺少实际工具类用法

- [ ] `bytes.Buffer`、`strings.Reader`、`io.Copy`
- [ ] `fmt.Fprintf` / `fmt.Fscanf` — 格式化读写
- [ ] `bufio.Scanner` — 逐行扫描大文件
- [ ] 练习：实现一个简单的数据流管道

### 补充模块：net/http 基础

> 第三阶段 HTTP/gRPC 的前置

- [ ] `http.Get` / `http.NewRequest` — 客户端
- [ ] `http.Handler` / `http.HandleFunc` — 服务端
- [ ] `http.Client` 超时设置
- [ ] 练习：简易 HTTP API 客户端

---

## 剩余路线

### 错误处理（部分已掌握）

- [x] `errors.Is`、`errors.As`、`%w` — 第一阶段
- [x] 自定义错误类型 — 练习19
- [ ] 练习：带 super时和重试机制的 HTTP 请求（待学 http 模块后做）

### 并发 🔥 分布式核心（待全部前置模块完成后）

- [ ] **2.8** Goroutine + Channel + select
- [ ] **2.9** Pipeline、Fan-out/Fan-in
- [ ] **2.10** sync.Mutex、WaitGroup、Once、atomic
- [ ] **2.11** Context 包

### 分布式专项

- [ ] 3.1 网络编程（TCP/HTTP）
- [ ] 3.2 gRPC（protobuf + Stream）
- [ ] 3.3 服务注册与发现（etcd/consul）
- [ ] 3.4 负载均衡
- [ ] 3.5 分布式共识（Raft）
- [ ] 3.6 消息队列（Pub/Sub）
- [ ] 3.7 可观测性（OpenTelemetry）

### 综合项目

- [ ] 4.1 分布式 KV 存储
- [ ] 4.2 任务调度系统

---

## 推荐资源

- **并发**: [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- **Raft**: [Raft 可视化](https://raft.github.io/)
- **分布式**: [MIT 6.824](https://pdos.csail.mit.edu/6.824/)
- **gRPC**: [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
