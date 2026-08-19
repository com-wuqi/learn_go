package practice

import (
	"context"
	"errors"
	"fmt"
	"net"

	"LearnGo/api/calc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ============================================================
// 第三阶段 3.2：gRPC Unary RPC
// ============================================================

// calculatorServer 实现 calc.CalculatorServer 接口。
// 提示：嵌入 calc.UnimplementedCalculatorServer 以获得向前兼容。
type calculatorServer struct {
	// TODO: 嵌入未实现类型
	calc.UnimplementedCalculatorServer
}

// Add 实现两个 int64 的加法。
// 提示：把 req.A + req.B 写进 AddReply.Result，错误返回 nil。
func (s *calculatorServer) Add(ctx context.Context, req *calc.AddRequest) (*calc.AddReply, error) {
	// TODO: 实现加法
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &calc.AddReply{Result: req.A + req.B}, nil
	}

}

// RunCalculatorDemo 起一个 Calculator server，再用客户端调用 Add。
// 完整流程参考 review/grpc_unary.go 的 RunUnaryDemo：
//  1. net.Listen 监听 127.0.0.1:0
//  2. grpc.NewServer + calc.RegisterCalculatorServer 注册
//  3. 后台 Serve，最后 Stop 优雅关闭
//  4. grpc.NewClient + insecure.NewCredentials 建立明文连接
//  5. calc.NewCalculatorClient + client.Add 发起调用并打印结果
//
// 实现本函数时需要补充 import：fmt、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorDemo() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer func(ln net.Listener) {
	//	err := ln.Close()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}(ln) // server.GracefulStop()
	server := grpc.NewServer()
	calc.RegisterCalculatorServer(server, &calculatorServer{})
	go func() {
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer server.GracefulStop()
	client, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(client *grpc.ClientConn) {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(client)
	ans, err := calc.NewCalculatorClient(client).Add(context.Background(), &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ans)

}
