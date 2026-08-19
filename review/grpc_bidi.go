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

// SayHelloBidi 是 Bidirectional-streaming 的实现：
// 服务端逐条 Recv，每收到一条就 Send 一条回复；客户端 CloseSend 后 Recv 到 io.EOF 即结束。
func (s *greeterServer) SayHelloBidi(stream demo.Greeter_SayHelloBidiServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&demo.HelloReply{Message: "Hi, " + req.GetName() + "!"}); err != nil {
			return err
		}
	}
}

// RunBidiDemo 演示 Bidirectional-streaming RPC：
// 客户端开一个 goroutine 收响应，主协程发请求；发完 CloseSend，等收方 io.EOF 结束。
func RunBidiDemo() {
	fmt.Println("\n=== gRPC Bidirectional-streaming 演示：Greeter.SayHelloBidi ===")

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := client.SayHelloBidi(ctx)
	if err != nil {
		fmt.Println("  发起流失败:", err)
		return
	}

	// 接收 goroutine：持续 Recv 直到 io.EOF，通知主协程。
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			reply, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Println("  Recv 失败:", err)
				return
			}
			fmt.Println("  客户端收到:", reply.GetMessage())
		}
	}()

	// 主协程逐条发送；全双工下发送和接收互不阻塞。
	for _, name := range []string{"Alice", "Bob", "Carol"} {
		if err := stream.Send(&demo.HelloRequest{Name: name}); err != nil {
			fmt.Println("  Send 失败:", err)
			return
		}
		fmt.Println("  客户端发送:", name)
	}

	// 关闭发送方向：服务端 Recv 会收到 io.EOF。
	if err := stream.CloseSend(); err != nil {
		fmt.Println("  CloseSend 失败:", err)
		return
	}
	<-done
}
