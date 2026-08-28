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

// loggingStreamServerInterceptor 包装整个流式 RPC：进入时记录，结束时打印耗时。
// 注意：流式拦截器每条 RPC（每条流）触发一次，不是每条消息。
func loggingStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, ss)
	fmt.Printf("  [server stream interceptor] %s 耗时=%s\n", info.FullMethod, time.Since(start).Round(time.Millisecond))
	return err
}

// loggingStreamClientInterceptor 包装客户端流：建立流时打印方法名。
func loggingStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream, err := streamer(ctx, desc, cc, method, opts...)
	fmt.Printf("  [client stream interceptor] %s\n", method)
	return stream, err
}

// RunStreamInterceptorDemo 演示流式拦截器：一次 Server-streaming 调用，两侧各打一条日志。
func RunStreamInterceptorDemo() {
	fmt.Println("\n=== gRPC 流式拦截器演示：Server-streaming 两侧日志 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer(grpc.ChainStreamInterceptor(loggingStreamServerInterceptor))
	demo.RegisterGreeterServer(srv, &greeterServer{})

	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.NewClient(ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(loggingStreamClientInterceptor),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := demo.NewGreeterClient(conn).SayHelloStream(ctx, &demo.HelloRequest{Name: "Codex", Count: 3})
	if err != nil {
		fmt.Println("  发起流失败:", err)
		return
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("  Recv 失败:", err)
			return
		}
		fmt.Println("  收到:", reply.GetMessage())
	}
}
