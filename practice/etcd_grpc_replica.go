package practice

import (
	"context"
	"fmt"
	"sync"

	"LearnGo/api/demo"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

// ============================================================
// 第三阶段 3.3-E3A：单副本注册式 gRPC 服务
//
// 一个副本必须同时管理三类资源：
//   - net.Listener / grpc.Server：提供实际 RPC
//   - etcd lease 与注册 key：告诉发现端“我在哪里”
//   - 心跳注销函数：副本下线时撤租约并清理 key
//
// 参考：review/grpc_unary.go + practice/etcd_register.go
// ============================================================

type RegisteredGreeter struct {
	Addr       string
	InstanceID string

	server           *grpc.Server
	stopRegistration func() error
	stopOnce         sync.Once
	stopErr          error
}

// labeledGreeter 的响应携带副本标签，便于集成测试确认请求落到哪个副本。
// TODO: 实现 SayHello，返回包含 label 和请求 Name 的稳定文本。
type labeledGreeter struct {
	demo.UnimplementedGreeterServer
	label string
}

func (g *labeledGreeter) SayHello(ctx context.Context, req *demo.HelloRequest) (*demo.HelloReply, error) {
	// TODO
	return nil, nil
}

// StartRegisteredGreeter 启动并注册一个 gRPC 副本。
//
// 契约：
//   - net.Listen("tcp", "127.0.0.1:0")，用实际 ln.Addr().String() 作为 etcd value
//   - grpc.NewServer + demo.RegisterGreeterServer + 后台 Serve
//   - RegisterServiceWithHeartbeat(ctx, cli, service, instanceID, addr, ttl)
//   - 注册失败时停止 server、关闭 listener，不泄漏资源
//   - 成功返回后服务可立即被调用
func StartRegisteredGreeter(ctx context.Context, cli *clientv3.Client, service, instanceID, label string, ttl int64) (*RegisteredGreeter, error) {
	// TODO
	return nil, fmt.Errorf("TODO: implement StartRegisteredGreeter")
}

// Stop 注销并停止副本；必须幂等，且重复调用返回第一次的错误。
//
// 建议顺序：
//  1. stopRegistration（停止心跳并 Revoke，等待完成）
//  2. server.Stop（停止监听并断开连接）
//  3. errors.Join 合并错误
func (r *RegisteredGreeter) Stop() error {
	// TODO
	return nil
}

// 自检：
//   - Start 后 etcd key 的 value 等于 r.Addr
//   - 直接 grpc.NewClient(r.Addr) 的 SayHello 能返回副本 label
//   - 第一次 Stop 后 etcd key 消失，再次 Stop 不重复 Revoke
//   - 注册失败或测试提前退出时不留下 listener/server/lease
// ------------------------------------------------------------
