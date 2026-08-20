package review

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// greeterServer 实现 demo.GreeterServer 接口。
// 必须嵌入 demo.UnimplementedGreeterServer，这是 gRPC 的向前兼容机制：
// 以后 Greeter 新增方法时，旧实现不会因为缺方法而编译失败。
type greeterServer struct {
	demo.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *demo.HelloRequest) (*demo.HelloReply, error) {
	// 演示超时：SleepMs > 0 时服务端故意延迟，期间响应 ctx 取消。
	if req.GetSleepMs() > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(req.GetSleepMs()) * time.Millisecond):
		}
	}
	// 演示 metadata：拦截器把 x-user-id 放进了 ctx。
	if user, ok := ctx.Value(userKey).(string); ok && user != "" {
		return &demo.HelloReply{Message: fmt.Sprintf("Hello, %s! (user=%s)", req.GetName(), user)}, nil
	}
	return &demo.HelloReply{Message: "Hello, " + req.GetName() + "!"}, nil
}

// RunUnaryDemo 演示一个完整的 Unary RPC：
// 起 server -> 建 client 连接 -> 调用 -> 优雅关闭。
func RunUnaryDemo() {
	fmt.Println("\n=== gRPC Unary RPC 演示：Greeter.SayHello ===")

	// 1. 监听一个随机的本机端口
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	// 2. 创建 gRPC server，并注册我们的实现
	srv := grpc.NewServer()
	demo.RegisterGreeterServer(srv, &greeterServer{})

	// 3. 在后台启动 Serve
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv.Stop()

	// 4. 建立客户端连接。演示环境没有 TLS，所以用 insecure。
	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	// 5. 用生成的客户端发起 RPC
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
