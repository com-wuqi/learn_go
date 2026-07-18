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
- `review/*.go` — 知识点讲解代码
>   - `practice/exercises.go` — 练习 1-9
>   - `practice/interfaces.go` — 练习 15-22
>   - `practice/io.go` — 练习 23
>   - `practice/structs.go` — 练习 24-30
>   - `practice/reflect.go` — 练习 31-33, 36
>   - `practice/kvstore_test.go` — 练习 34-35
>   - `practice/time.go` — 练习 37-38
>   - `practice/concurrency.go` — 练习 10-14 待学

---

## 🎯 当前进度

| 阶段       | 状态         |
|----------|------------|
| 第一阶段     | ✅          |
| **第二阶段** | **🔄 复习中** |
| 第三阶段     | 待开始        |
| 第四阶段     | 待开始        |

### 🔥 复习：二阶段前半程总结

- [ ] 复习题1: 完整 KVStore 测试套件（含 JSON 文件持久化验证）
- [ ] 复习题2: 管道 — 读文件 → Plugin 转换 → 写文件
- [ ] 复习题3: ConfigLoader — 反射 + JSON + default tag 填充

### 🔥 io.Reader/Writer 实战 ✅ 已完成

- [x] 阅读 `review/io_reader.go`
- [x] 练习39: CountLines
- [x] 练习40: FilterLines
- [x] 练习41: ConcatReaders

### 🔥 net/http 基础

- [ ] 阅读 `review/http.go` — 跑 `go run ./cmd/phase2/`
- [ ] 练习42: HTTP GET 请求
- [ ] 练习43: 简易 HTTP 服务端
- [ ] 练习44: 带超时的 HTTP 客户端

### 下一步

- [ ] net/http 基础
- [ ] 并发（goroutine → channel → select → sync → context）

---

## 剩余路线

### 并发 🔥 分布式核心

- [ ] Goroutine + Channel + select
- [ ] Pipeline、Fan-out/Fan-in
- [ ] sync.Mutex、WaitGroup、Once、atomic
- [ ] Context 包

### 分布式专项

- [ ] 3.1 网络编程（TCP/HTTP）
- [ ] 3.2 gRPC
- [ ] 3.3 服务注册与发现
- [ ] 3.4 负载均衡
- [ ] 3.5 分布式共识（Raft）
- [ ] 3.6 消息队列
- [ ] 3.7 可观测性

### 综合项目

- [ ] 4.1 分布式 KV 存储
- [ ] 4.2 任务调度系统

## 推荐资源

- **并发**: [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- **Raft**: [Raft 可视化](https://raft.github.io/)
- **分布式**: [MIT 6.824](https://pdos.csail.mit.edu/6.824/)
