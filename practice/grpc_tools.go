package practice

import (
	"LearnGo/api/calc"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/stats"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ============================================================
// 第三阶段 3.2-S 批3：流式拦截器 / 压缩与限长 / health / 反射
// ============================================================

// loggingStreamServerInterceptor 服务端流式拦截器：打印方法名与耗时。
// 签名：func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
// 注意：流式拦截器每条 RPC（每条流）触发一次，不是每条消息。
func loggingStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// TODO
	start := time.Now()
	err := handler(srv, ss)
	duration := time.Since(start).Round(time.Millisecond)
	fmt.Printf("  [server stream interceptor] %s 耗时=%s\n", info.FullMethod, duration)
	return err
}

// loggingStreamClientInterceptor 客户端流式拦截器：建立流时打印方法名。
// 签名：func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error)
func loggingStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// TODO
	var stream grpc.ClientStream
	var err error
	start := time.Now()
	stream, err = streamer(ctx, desc, cc, method, opts...)
	duration := time.Since(start).Round(time.Millisecond)
	fmt.Printf("  [client stream interceptor] %s 耗时=%s\n", method, duration)
	return stream, err
}

// RunStreamInterceptorPractice 起带流式拦截器的 server/client，调用 EchoSum 或 ListNumbers，观察两侧日志。
// 完整流程参考 review/grpc_stream_interceptor.go。
func RunStreamInterceptorPractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	server := grpc.NewServer(grpc.StreamInterceptor(loggingStreamServerInterceptor))
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
		grpc.WithStreamInterceptor(loggingStreamClientInterceptor),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	echoSum, err := client.EchoSum(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			resp, err := echoSum.Recv()
			// If io.EOF is returned, the stream has terminated with an OK status.
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(resp.GetResult())

		}
	}()
	for i := 0; i < 100; i++ {
		err = echoSum.Send(&calc.AddRequest{
			A: int64(i),
			B: int64(2 * i),
		})
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	err = echoSum.CloseSend()
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()

}

type sizeStatsHandler struct{}

func (s *sizeStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}
func (s *sizeStatsHandler) HandleRPC(ctx context.Context, r stats.RPCStats) {
	if p, ok := r.(*stats.InPayload); ok {
		ratio := float64(p.CompressedLength) / float64(p.Length) * 100
		fmt.Printf("  [server stats] 收到 %d 字节，压缩后 %d 字节（%.1f%%）\n", p.Length, p.CompressedLength, ratio)
	}
}
func (s *sizeStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}
func (s *sizeStatsHandler) HandleConn(ctx context.Context, c stats.ConnStats) {}

// RunSizesPractice 起带 MaxRecvMsgSize 限制的 server，分别用无压缩/gzip 调用，观察 stats 输出；
// 再发一个超限消息，观察 codes.ResourceExhausted。
// 完整流程参考 review/grpc_sizes.go。
func RunSizesPractice() {
	// TODO
	// 服务端会尝试用和请求相同的编码回复，导入包是支持解码请求和发回
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(1<<10), grpc.StatsHandler(&sizeStatsHandler{}))
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
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.Add(ctx, &calc.AddRequest{A: 1, B: 2, Payload: "Hello"})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.GetResult())
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	resp2, err := client.Add(ctx2, &calc.AddRequest{A: 1, B: 2, Payload: strings.Repeat("x", 1<<10)})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.GetResult())
	}
}

// RunHealthPractice 注册 health server 并查询状态（SERVING → 改 NOT_SERVING → 停服后 Unavailable）。
// 完整流程参考 review/grpc_health.go。
func RunHealthPractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	healthServer := health.NewServer()
	healthServer.SetServingStatus("Calculator.Add", healthpb.HealthCheckResponse_SERVING)
	server := grpc.NewServer()
	calc.RegisterCalculatorServer(server, &calculatorServer{})
	healthpb.RegisterHealthServer(server, healthServer)
	go func() {
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
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
	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	healthCheckRequest := &healthpb.HealthCheckRequest{Service: "Calculator.Add"}
	resp, err := client.Check(ctx, healthCheckRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.GetStatus())
	}
	healthServer.SetServingStatus("Calculator.Add", healthpb.HealthCheckResponse_NOT_SERVING)
	resp2, err := client.Check(ctx, healthCheckRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.GetStatus())
	}

	server.GracefulStop()
	resp3, err := client.Check(ctx, healthCheckRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp3.GetStatus())
	}

}

// RunReflectionPractice 开反射并用反射客户端列出服务与方法。
// 完整流程参考 review/grpc_reflection.go。
func RunReflectionPractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := grpc.NewServer()
	calc.RegisterCalculatorServer(srv, &calculatorServer{})
	reflection.Register(srv)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	clientConn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := grpc_reflection_v1.NewServerReflectionClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.ServerReflectionInfo(ctx) // 反射流
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{ListServices: ""},
	}); err != nil {
		fmt.Println(err)
		return
	}
	resp, err := stream.Recv()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, svc := range resp.GetListServicesResponse().GetService() {
		fmt.Println("svc:", svc.GetName())
	}

	if err := stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: "calc.Calculator"},
	}); err != nil {
		fmt.Println("  Send 失败:", err)
		return
	}
	resp, err = stream.Recv()
	if err != nil {
		fmt.Println("  Recv 失败:", err)
		return
	}
	fdBytes := resp.GetFileDescriptorResponse().GetFileDescriptorProto()
	if len(fdBytes) == 0 {
		return
	}
	fd := &descriptorpb.FileDescriptorProto{}
	if err := proto.Unmarshal(fdBytes[0], fd); err != nil {
		fmt.Println("  解析文件描述符失败:", err)
		return
	}
	for _, svc := range fd.GetService() {
		for _, m := range svc.GetMethod() {
			fmt.Printf("%s (%s) -> (%s)\n", m.GetName(), m.GetInputType(), m.GetOutputType())
		}
	}
}
