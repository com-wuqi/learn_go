package review

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SayHelloClientStream 是 Client-streaming 的实现：
// 持续 Recv 收集所有名字，最后 SendAndClose 返回一条汇总响应。
func (s *greeterServer) SayHelloClientStream(stream demo.Greeter_SayHelloClientStreamServer) error {
	var names []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&demo.HelloReply{Message: "Hello, " + strings.Join(names, ", ") + "!"})
		}
		if err != nil {
			return err
		}
		names = append(names, req.GetName())
	}
}

// RunClientStreamDemo 演示 Client-streaming RPC：
// 客户端循环 Send 多个请求，最后 CloseAndRecv 拿一个响应。
func RunClientStreamDemo() {
	fmt.Println("\n=== gRPC Client-streaming 演示：Greeter.SayHelloClientStream ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer()
	demo.RegisterGreeterServer(srv, &greeterServer{})

	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	client := demo.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := client.SayHelloClientStream(ctx)
	if err != nil {
		fmt.Println("  发起流失败:", err)
		return
	}

	for _, name := range []string{"Alice", "Bob", "Carol"} {
		if err := stream.Send(&demo.HelloRequest{Name: name}); err != nil {
			fmt.Println("  Send 失败:", err)
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("  CloseAndRecv 失败:", err)
		return
	}
	fmt.Println("  客户端收到:", resp.GetMessage())
}
