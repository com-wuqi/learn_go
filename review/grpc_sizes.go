package review

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip" // 注册 gzip 压缩器（同时提供 gzip.Name）
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// sizeStatsHandler 统计服务端收到的消息大小（原始 vs 压缩后）。
type sizeStatsHandler struct{}

func (s *sizeStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}
func (s *sizeStatsHandler) HandleRPC(ctx context.Context, r stats.RPCStats) {
	if p, ok := r.(*stats.InPayload); ok {
		ratio := float64(p.CompressedLength) / float64(p.Length) * 100
		fmt.Printf("  [server stats] 收到 %d 字节，压缩后 %d 字节（%.1f%%）\n", p.Length, p.CompressedLength, ratio)
	}
}
func (s *sizeStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}
func (s *sizeStatsHandler) HandleConn(ctx context.Context, c stats.ConnStats) {}

// RunSizesDemo 演示 gzip 压缩与消息大小限制。
func RunSizesDemo() {
	fmt.Println("\n=== gRPC 压缩与消息大小限制演示 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}
	srv := grpc.NewServer(grpc.StatsHandler(&sizeStatsHandler{}))
	demo.RegisterGreeterServer(srv, &greeterServer{})
	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addr := ln.Addr().String()
	bigName := strings.Repeat("x", 10240)

	// 1. 不压缩
	connPlain, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	resp1, err := demo.NewGreeterClient(connPlain).SayHello(ctx, &demo.HelloRequest{Name: bigName})
	if err != nil {
		fmt.Println("  1) 无压缩调用失败:", err)
	} else {
		fmt.Println("  1) 无压缩:", resp1.GetMessage()[:12]+"...")
	}
	connPlain.Close()

	// 2. gzip 压缩（通过默认调用选项）
	connZip, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	resp2, err := demo.NewGreeterClient(connZip).SayHello(ctx, &demo.HelloRequest{Name: bigName})
	if err != nil {
		fmt.Println("  2) gzip 调用失败:", err)
	} else {
		fmt.Println("  2) gzip 压缩:", resp2.GetMessage()[:12]+"...")
	}
	connZip.Close()

	// 3. 服务端限制接收 1KB，10KB 消息被拒
	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}
	srv2 := grpc.NewServer(grpc.MaxRecvMsgSize(1 << 10))
	demo.RegisterGreeterServer(srv2, &greeterServer{})
	go func() {
		if err := srv2.Serve(ln2); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv2.Stop()

	conn3, err := grpc.NewClient(ln2.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn3.Close()
	_, err = demo.NewGreeterClient(conn3).SayHello(ctx, &demo.HelloRequest{Name: bigName})
	fmt.Println("  3) 超过 MaxRecvMsgSize(1KB):", status.Code(err))
}
