package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ============================================================
// 第三阶段 3.2：gRPC metadata / 拦截器（可选进阶）
// ============================================================

// loggingUnaryServerInterceptor 服务端 Unary 拦截器：
// 打印方法名和耗时；若 incoming metadata 带 x-user-id，也打印出来。
// 签名：func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
// 提示：metadata.FromIncomingContext(ctx) 取 metadata；handler(ctx, req) 调用链上的下一个处理器。
func loggingUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	begin := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(begin).Round(time.Millisecond)
	fmt.Printf("info.FullMethod: %s, duration: %s\n", info.FullMethod, duration)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val := md.Get("x-user-id"); len(val) > 0 {
			fmt.Printf("[x-user-id]: %s\n", val[0])
		}
	}
	return resp, err
}

// loggingUnaryClientInterceptor 客户端 Unary 拦截器：打印方法名和耗时。
// 签名：func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
// 提示：invoker(ctx, method, req, reply, cc, opts...) 发起真实调用。
func loggingUnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// TODO
	begin := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Since(begin).Round(time.Millisecond)
	fmt.Printf("method: %s, duration: %s\n", method, duration)
	return err

}

// RunCalculatorAdvancedDemo 起 server（带服务端拦截器）、建 client（带客户端拦截器），
// 用 metadata.AppendToOutgoingContext 附加 x-user-id 后调用 Add，观察两侧日志。
// 完整流程参考 review/grpc_advanced.go 的 RunMetadataDemo / RunInterceptorDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、
// google.golang.org/grpc/metadata、google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorAdvancedDemo() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingUnaryServerInterceptor))
	calc.RegisterCalculatorServer(server, &calculatorServer{})
	go func() {
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer server.GracefulStop()
	conn, err := grpc.NewClient(
		ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(loggingUnaryClientInterceptor),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	client := calc.NewCalculatorClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-user-id", "1")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	addR, err := client.Add(ctx, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println("  客户端收到错误:", err)
		fmt.Println("  状态码:", status.Code(err))
		return
	}
	fmt.Println(addR.GetResult())

}
