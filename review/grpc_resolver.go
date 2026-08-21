package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// staticResolverBuilder 返回固定地址列表的 resolver，scheme 为 "demo"。
// 自定义 resolver 是实现服务发现（如 etcd）的入口：Build 里拿到地址后喂给 ClientConn。
type staticResolverBuilder struct {
	addrs []string
}

func (b *staticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0, len(b.addrs))
	for _, a := range b.addrs {
		addrs = append(addrs, resolver.Address{Addr: a})
	}
	cc.UpdateState(resolver.State{Addresses: addrs})
	return &staticResolver{cc: cc}, nil
}

func (b *staticResolverBuilder) Scheme() string { return "demo" }

type staticResolver struct {
	cc resolver.ClientConn
}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *staticResolver) Close()                                {}

// runBackend 起一个带标签的 greeter server，返回监听地址和停止函数。
func runBackend(label string) (string, func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		fmt.Printf("  [backend %s] 收到 %s\n", label, info.FullMethod)
		return handler(ctx, req)
	}))
	demo.RegisterGreeterServer(srv, &greeterServer{})
	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	return ln.Addr().String(), srv.Stop
}

// RunResolverDemo 演示自定义 resolver + round_robin：一个客户端把请求轮询到两个后端。
func RunResolverDemo() {
	fmt.Println("\n=== gRPC resolver + 负载均衡演示：round_robin 分发到两个后端 ===")

	addrA, stopA := runBackend("A")
	defer stopA()
	addrB, stopB := runBackend("B")
	defer stopB()

	builder := &staticResolverBuilder{addrs: []string{addrA, addrB}}
	conn, err := grpc.NewClient(
		"demo:///lb",
		grpc.WithResolvers(builder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	client := demo.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 6; i++ {
		resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
		if err != nil {
			fmt.Println("  调用失败:", err)
			return
		}
		fmt.Printf("  第 %d 次调用: %s\n", i, resp.GetMessage())
		time.Sleep(100 * time.Millisecond)
	}
}
