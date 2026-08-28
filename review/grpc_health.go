package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// RunHealthDemo 演示标准健康检查协议 grpc_health_v1。
func RunHealthDemo() {
	fmt.Println("\n=== gRPC health 检查演示：服务状态查询与变更 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	hs := health.NewServer()
	hs.SetServingStatus("hello.Greeter", healthpb.HealthCheckResponse_SERVING)

	srv := grpc.NewServer()
	demo.RegisterGreeterServer(srv, &greeterServer{})
	healthpb.RegisterHealthServer(srv, hs)
	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()

	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: "hello.Greeter"})
	fmt.Println("  1) 服务状态:", resp.GetStatus(), err)

	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "unknown.Service"})
	fmt.Println("  2) 未知服务:", status.Code(err))

	hs.SetServingStatus("hello.Greeter", healthpb.HealthCheckResponse_NOT_SERVING)
	resp, _ = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "hello.Greeter"})
	fmt.Println("  3) 改为 NOT_SERVING:", resp.GetStatus())

	srv.Stop()
	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "hello.Greeter"})
	fmt.Println("  4) 服务停止后:", status.Code(err))
}
