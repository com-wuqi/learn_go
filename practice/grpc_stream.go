package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ============================================================
// 第三阶段 3.2：gRPC Server-streaming RPC
// ============================================================

// ListNumbers 实现服务端流：把 [from, to] 区间内的整数逐个 Send 出去。
// 提示：stream 类型是 calc.Calculator_ListNumbersServer，
//
//	用 stream.Send(&calc.NumberReply{Value: i}) 发送。
func (s *calculatorServer) ListNumbers(req *calc.RangeRequest, stream calc.Calculator_ListNumbersServer) error {
	// TODO
	ctx := stream.Context()
	for i := req.GetFrom(); i <= req.GetTo(); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			{
				err := stream.Send(&calc.NumberReply{
					Value: i,
				})
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

// RunCalculatorStreamDemo 起 server、建 client、调用 ListNumbers 并逐个 Recv。
// 完整流程参考 review/grpc_stream.go 的 RunServerStreamDemo。
//
// 实现本函数时需要补充 import：context、fmt、io、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorStreamDemo() {
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
	stream, err := calc.NewCalculatorClient(client).ListNumbers(context.Background(), &calc.RangeRequest{From: 1, To: 20})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		rasp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Recv 失败:", err)
			return
		}
		fmt.Println(rasp.GetValue())
	}
}
