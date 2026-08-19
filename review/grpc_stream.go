package review

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SayHelloStream 是 Server-streaming 的实现：
// 收到一个请求，按 req.count 条数持续往 stream 里 Send 多条响应。
func (s *greeterServer) SayHelloStream(req *demo.HelloRequest, stream demo.Greeter_SayHelloStreamServer) error {
	for i := 0; i < int(req.GetCount()); i++ {
		if err := stream.Send(&demo.HelloReply{Message: fmt.Sprintf("Hello #%d, %s!", i+1, req.GetName())}); err != nil {
			return err
		}
	}
	return nil
}

// RunServerStreamDemo 演示 Server-streaming RPC：
// 客户端发一个请求，然后循环 Recv 直到 io.EOF。
func RunServerStreamDemo() {
	fmt.Println("\n=== gRPC Server-streaming 演示：Greeter.SayHelloStream ===")

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

	stream, err := client.SayHelloStream(ctx, &demo.HelloRequest{Name: "Codex", Count: 3})
	if err != nil {
		fmt.Println("  发起流失败:", err)
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("  流结束")
			break
		}
		if err != nil {
			fmt.Println("  Recv 失败:", err)
			return
		}
		fmt.Println("  收到:", resp.GetMessage())
	}
}
