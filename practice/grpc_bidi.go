package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ============================================================
// 第三阶段 3.2：gRPC Bidirectional-streaming RPC
// ============================================================

// EchoSum 实现双向流：循环 Recv，把每个 AddRequest 的 a+b 逐条回发。
// 提示：
//   - stream 类型是 calc.Calculator_EchoSumServer
//   - 服务端 Recv 到 io.EOF（客户端 CloseSend）时 return nil
//   - 每条消息用 stream.Send(&calc.AddReply{Result: req.A + req.B}) 回发
func (s *calculatorServer) EchoSum(stream calc.Calculator_EchoSumServer) error {
	// TODO
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			{
				req, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					return nil
				}
				if err != nil {
					return err
				}
				err = stream.Send(&calc.AddReply{Result: req.GetA() + req.GetB()})
				if err != nil {
					return err
				}
			}
		}
	}
}

// RunCalculatorBidiDemo 起 server、建 client、边发边收，最后 CloseSend。
// 完整流程参考 review/grpc_bidi.go 的 RunBidiDemo。
//
// 实现本函数时需要补充 import：context、fmt、io、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorBidiDemo() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcServer := grpc.NewServer()
	calc.RegisterCalculatorServer(grpcServer, &calculatorServer{})
	go func() {
		err := grpcServer.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer grpcServer.GracefulStop()
	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	client := calc.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stream, err := client.EchoSum(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	group := errgroup.Group{}
	group.Go(func() error {
		for {
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				//fmt.Println(err)
				return err
			}
			fmt.Println(req.GetResult())
		}
	})
	for _, d := range []int{1, 2, 3, 4, 5} {
		err := stream.Send(&calc.AddRequest{
			A: int64(d),
			B: int64(0),
		})
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	err = stream.CloseSend()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := group.Wait(); err != nil {
		fmt.Println(err)
		return
	}

}
