package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ctxKey 用于把 metadata 解析出的用户信息放进 context。
type ctxKey string

// context.WithValue 进程内沿调用链传数据
// metadata.AppendToOutgoingContext 把 KV 打包成 HTTP/2 头发给对端

const userKey ctxKey = "user"

// metadataUnaryServerInterceptor 从 incoming metadata 读取 x-user-id，注入 ctx。
func metadataUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	user := "anonymous"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-user-id"); len(vals) > 0 {
			user = vals[0]
		}
	}
	ctx = context.WithValue(ctx, userKey, user)
	return handler(ctx, req)
}

// loggingUnaryServerInterceptor 打印每个 Unary 调用的方法和耗时。
func loggingUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	fmt.Printf("  [server interceptor] %s 耗时=%s\n", info.FullMethod, time.Since(start).Round(time.Millisecond))
	return resp, err
}

// unaryClientInterceptor 打印客户端发起的每个 Unary 调用的方法和耗时。
func unaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("  [client interceptor] %s 耗时=%s err=%v\n", method, time.Since(start).Round(time.Millisecond), err)
	return err
}

// RunMetadataDemo 演示 metadata：客户端用 metadata.AppendToOutgoingContext 附加键值对，
// 服务端拦截器用 metadata.FromIncomingContext 取出并注入 ctx，SayHello 读取后回显。
func RunMetadataDemo() {
	fmt.Println("\n=== gRPC metadata 演示：客户端附 x-user-id，服务端读取并回显 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(metadataUnaryServerInterceptor, loggingUnaryServerInterceptor))
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

	// 把键值对附加到 outgoing metadata。
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-user-id", "alice")
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Alice"})
	if err != nil {
		fmt.Println("  SayHello 失败:", err)
		return
	}
	fmt.Println("  客户端收到:", resp.GetMessage())
}

// RunTimeoutDemo 演示超时传播：客户端给 1s deadline，服务端故意 sleep 3s。
// 客户端在本地先超时，服务端也会通过 ctx.Done 感知取消。
func RunTimeoutDemo() {
	fmt.Println("\n=== gRPC 超时演示：客户端 1s deadline，服务端 sleep 3s ===")

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Slow", SleepMs: 3000})
	fmt.Printf("  客户端实际等待: %s\n", time.Since(start).Round(time.Millisecond))
	if err != nil {
		fmt.Println("  客户端收到错误:", err)
		fmt.Println("  状态码:", status.Code(err))
		return
	}
	fmt.Println("  客户端收到:", resp.GetMessage())
}

// RunInterceptorDemo 演示拦截器：服务端和客户端各挂一个日志拦截器，观察两侧输出。
func RunInterceptorDemo() {
	fmt.Println("\n=== gRPC 拦截器演示：两侧打印方法名与耗时 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingUnaryServerInterceptor))
	demo.RegisterGreeterServer(srv, &greeterServer{})

	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.NewClient(
		ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(unaryClientInterceptor),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	client := demo.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		fmt.Println("  SayHello 失败:", err)
		return
	}
	fmt.Println("  客户端收到:", resp.GetMessage())
}
