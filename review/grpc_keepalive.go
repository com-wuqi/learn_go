package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// RunKeepaliveDemo 演示 keepalive 的两个层面：
// 1) 服务端 KeepaliveParams.MaxConnectionIdle：空闲连接超时回收；
// 2) 客户端 WithKeepaliveParams：心跳参数（grpc-go 对客户端 ping 间隔有 10s 硬性下限）。
// 关于 too_many_pings：服务端 EnforcementPolicy.MinTime 限制客户端 ping 频率，
// 违反时服务端发 GOAWAY(ENHANCE_YOUR_CALM, "too_many_pings") 断开连接。
func RunKeepaliveDemo() {
	fmt.Println()
	fmt.Println("=== gRPC keepalive 演示：服务端空闲回收与心跳参数 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer(
		// 服务端主动心跳：每 1s 探测一次对端（Time 有 1s 下限）。
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              1 * time.Second,
			Timeout:           1 * time.Second,
			MaxConnectionIdle: 3 * time.Second, // 空闲 3s 后关闭连接
		}),
		// 服务端 enforcement：限制客户端 ping 频率与无流 ping。
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	demo.RegisterGreeterServer(srv, &greeterServer{})
	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 客户端 keepalive：Time=10s 是 grpc-go 允许的最小间隔（更小会被抬到 10s），
	// 本演示窗口内看不到客户端心跳，重点是服务端 MaxConnectionIdle。
	conn, err := grpc.NewClient(ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             2 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	go watchConnState(conn, ctx, "conn")
	conn.Connect()

	client := demo.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		fmt.Println("  第一次 RPC 失败:", err)
		return
	}
	fmt.Println("  第一次 RPC:", resp.GetMessage())

	fmt.Println("  空闲等待：3s 后服务端应关闭空闲连接，观察状态变化...")
	time.Sleep(5 * time.Second)

	resp, err = client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		fmt.Println("  重连后 RPC 失败:", err)
		return
	}
	fmt.Println("  重连后 RPC:", resp.GetMessage())
}
