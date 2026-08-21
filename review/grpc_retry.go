package review

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// flakyServerInterceptor 前两次 SayHello 返回 Unavailable，之后放行。
type flakyServerInterceptor struct {
	mu    sync.Mutex
	count int
}

func (f *flakyServerInterceptor) unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == "/hello.Greeter/SayHello" {
		f.mu.Lock()
		f.count++
		attempt := f.count
		f.mu.Unlock()
		if attempt <= 2 {
			fmt.Printf("  [server] 第 %d 次调用：返回 Unavailable\n", attempt)
			return nil, status.Error(codes.Unavailable, "temporarily unavailable")
		}
		fmt.Printf("  [server] 第 %d 次调用：放行\n", attempt)
	}
	return handler(ctx, req)
}

// retryServiceConfig 给 hello.Greeter 的所有方法配置重试策略：
// 最多 4 次尝试，只重试 UNAVAILABLE，backoff 从 0.1s 起、倍增、上限 1s。
const retryServiceConfig = `{
  "methodConfig": [{
    "name": [{"service": "hello.Greeter"}],
    "retryPolicy": {
      "maxAttempts": 4,
      "initialBackoff": "0.1s",
      "maxBackoff": "1s",
      "backoffMultiplier": 2,
      "retryableStatusCodes": ["UNAVAILABLE"]
    }
  }]
}`

// RunRetryDemo 演示 gRPC 内置重试：先不带重试配置（失败），再带配置（自动重试成功）。
func RunRetryDemo() {
	fmt.Println("\n=== gRPC 重试演示：Unavailable 自动重试 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	flaky := &flakyServerInterceptor{}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(flaky.unary))
	demo.RegisterGreeterServer(srv, &greeterServer{})

	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	addr := ln.Addr().String()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 不带重试：第一次失败就直接返回。
	conn1, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	_, err = demo.NewGreeterClient(conn1).SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	fmt.Println("  无重试配置:", status.Code(err))
	conn1.Close()

	// 2. 带重试配置：前两次失败后自动重试，第三次成功。
	conn2, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryServiceConfig),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn2.Close()

	resp, err := demo.NewGreeterClient(conn2).SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		fmt.Println("  带重试仍失败:", err)
		return
	}
	fmt.Println("  带重试配置:", resp.GetMessage())
}
