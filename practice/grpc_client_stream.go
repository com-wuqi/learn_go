package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ============================================================
// 第三阶段 3.2：gRPC Client-streaming RPC
// ============================================================

// Sum 实现客户端流：循环 Recv 累加 value，最后 SendAndClose 一个 AddReply。
// 提示：stream 类型是 calc.Calculator_SumServer，
//
//	io.EOF 时用 stream.SendAndClose(&calc.NumberReply{Value: sums}) 结束。
func (s *calculatorServer) Sum(stream calc.Calculator_SumServer) error {
	// TODO
	sums := int64(0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&calc.NumberReply{
				Value: sums,
			})
		}
		if err != nil {
			return err
		}
		sums += req.GetA() + req.GetB()
	}
}

// RunCalculatorClientStreamDemo 起 server、建 client、循环 Send 多个数字，最后 CloseAndRecv。
// 完整流程参考 review/grpc_client_stream.go 的 RunClientStreamDemo。
//
// 实现本函数时需要补充 import：context、fmt、io、net、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorClientStreamDemo() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer ln.Close()
	server := grpc.NewServer()
	calc.RegisterCalculatorServer(server, &calculatorServer{})
	go func() {
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer server.GracefulStop()
	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stream, err := client.Sum(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, d := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		err := stream.Send(&calc.AddRequest{
			A: int64(d),
			B: int64(0),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.GetValue())
}
