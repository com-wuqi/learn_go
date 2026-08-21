package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// statusDemoServerInterceptor 根据请求内容返回不同的 status 错误，演示错误码分类。
func statusDemoServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	hr, ok := req.(*demo.HelloRequest)
	if !ok {
		return handler(ctx, req)
	}
	switch hr.GetName() {
	case "missing":
		return nil, status.Error(codes.NotFound, "user not found")
	case "forbidden":
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	if hr.GetAge() < 0 {
		st := status.New(codes.InvalidArgument, "age must be >= 0")
		st, _ = st.WithDetails(&demo.HelloReply{Message: "hint: age must be non-negative"})
		return nil, st.Err()
	}
	return handler(ctx, req)
}

// RunStatusDemo 演示客户端如何解析 gRPC 错误：code、message、details。
func RunStatusDemo() {
	fmt.Println("\n=== gRPC status 错误演示：错误码分类与解析 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(statusDemoServerInterceptor))
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

	cases := []*demo.HelloRequest{
		{Name: "Alice", Age: 18},
		{Name: "missing"},
		{Name: "forbidden"},
		{Name: "Bob", Age: -1},
	}
	for _, req := range cases {
		resp, err := client.SayHello(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			fmt.Printf("  code=%-16s msg=%q\n", st.Code(), st.Message())
			for _, d := range st.Details() {
				if h, ok := d.(*demo.HelloReply); ok {
					fmt.Printf("    detail: %s\n", h.GetMessage())
				}
			}
			continue
		}
		fmt.Println("  成功:", resp.GetMessage())
	}
}
