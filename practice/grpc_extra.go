package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ============================================================
// 第三阶段 3.2-S：gRPC 补充学习（status / 重试）
// ============================================================

// calcCheckServerInterceptor 服务端拦截器：当 AddRequest.A+B > 100 时返回
// status.Error(codes.InvalidArgument, "sum too large")，否则交给 handler。
// 进阶：用 status.New + WithDetails 附加一个 calc.AddReply 作为错误详情，再 st.Err()。
func calcCheckServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	addRequest, ok := req.(*calc.AddRequest)
	if !ok {
		return handler(ctx, req)
	}
	reqA := addRequest.GetA()
	reqB := addRequest.GetB()
	if reqA+reqB > 100 {
		st := status.New(codes.InvalidArgument, "sum too large")
		st, err := st.WithDetails(&calc.AddReply{Result: -1})
		if err != nil {
			return nil, err
		}
		return nil, st.Err()
	}
	return handler(ctx, req)
}

// RunStatusPractice 起带 calcCheckServerInterceptor 的 server，客户端分别用
// A+B=3（成功）和 A+B=101（失败）调用 Add；失败时用 status.FromError 打印 code 和 message。
// 完整流程参考 review/grpc_status.go 的 RunStatusDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure、
// google.golang.org/grpc/codes、google.golang.org/grpc/status。
func RunStatusPractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(calcCheckServerInterceptor))
	calc.RegisterCalculatorServer(grpcServer, &calculatorServer{})
	go func() {
		err := grpcServer.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
			return
		}
	}()
	defer grpcServer.GracefulStop()
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientReqs := []*calc.AddRequest{{A: 1, B: 2}, {A: 100, B: 1}}
	for _, clientReq := range clientReqs {
		clientResp, err := client.Add(ctx, clientReq)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				fmt.Println(err)
				continue
			}
			fmt.Printf("  code=%-16s msg=%q\n", st.Code(), st.Message())
			continue
		}
		fmt.Println(clientResp.GetResult())
	}

}

// flakyCalcServerInterceptor 前 2 次 Add 调用返回 Unavailable，之后放行。
// 提示：用 struct 字段 + sync.Mutex 计数，注意只对 "/calc.Calculator/Add" 生效；
// 别用包级裸变量，并发调用会 data race。
type flakyCalcServerInterceptor struct {
	mu    sync.Mutex
	count int
}

func (f *flakyCalcServerInterceptor) unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	if info.FullMethod == "/calc.Calculator/Add" {
		f.mu.Lock()
		f.count++
		count := f.count
		f.mu.Unlock()
		if count <= 2 {
			st := status.New(codes.Unavailable, "Server Unavailable")
			st, err := st.WithDetails(&calc.AddReply{Result: -1})
			if err != nil {
				return nil, err
			}
			return nil, st.Err()
		}
	}
	return handler(ctx, req)
}

// RunRetryPractice 起带 flakyCalcServerInterceptor 的 server：
// 1) 不带重试配置调 Add，观察 Unavailable 失败；
// 2) 带 retryPolicy（maxAttempts=4、retryableStatusCodes=["UNAVAILABLE"]）再调，观察自动重试成功。
// 完整流程参考 review/grpc_retry.go 的 RunRetryDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、sync、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure、
// google.golang.org/grpc/codes、google.golang.org/grpc/status。
const retryServiceConfig = `{
  "methodConfig": [{
    "name": [{"service": "calc.Calculator"}],
    "retryPolicy": {
      "maxAttempts": 4,
      "initialBackoff": "0.1s",
      "maxBackoff": "1s",
      "backoffMultiplier": 2,
      "retryableStatusCodes": ["UNAVAILABLE"]
    }
  }]
}`

func RunRetryPractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	flaky := &flakyCalcServerInterceptor{}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(flaky.unary))
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
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	client := calc.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientReqs := &calc.AddRequest{A: 1, B: 2}
	res, err := client.Add(ctx, clientReqs)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Println(err)
		} else {
			fmt.Printf("  code=%-16s msg=%q\n", st.Code(), st.Message())
		}
	} else {
		fmt.Println(res.GetResult())
	}
	//
	conn2, err := grpc.NewClient(
		ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryServiceConfig),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn2 *grpc.ClientConn) {
		err := conn2.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn2)
	client2 := calc.NewCalculatorClient(conn2)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	res2, err := client2.Add(ctx2, clientReqs)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Println(err)
		} else {
			fmt.Printf("  code=%-16s msg=%q\n", st.Code(), st.Message())
		}
	} else {
		fmt.Println(res2.GetResult())
	}

}
