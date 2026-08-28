package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// RunReflectionDemo 演示反射协议：客户端动态列出服务与方法（grpcurl 的底层机制）。
func RunReflectionDemo() {
	fmt.Println("\n=== gRPC 反射演示：动态列出服务与方法 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}
	srv := grpc.NewServer()
	demo.RegisterGreeterServer(srv, &greeterServer{})
	reflection.Register(srv)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client := reflectionpb.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		fmt.Println("  发起反射流失败:", err)
		return
	}

	// 1. 列出所有服务
	if err := stream.Send(&reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{ListServices: ""},
	}); err != nil {
		fmt.Println("  Send 失败:", err)
		return
	}
	resp, err := stream.Recv()
	if err != nil {
		fmt.Println("  Recv 失败:", err)
		return
	}
	for _, svc := range resp.GetListServicesResponse().GetService() {
		fmt.Println("  服务:", svc.GetName())
	}

	// 2. 用符号查文件描述符，解析出 hello.Greeter 的方法列表
	if err := stream.Send(&reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: "hello.Greeter"},
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
		fmt.Println("  没有拿到文件描述符")
		return
	}
	fd := &descriptorpb.FileDescriptorProto{}
	if err := proto.Unmarshal(fdBytes[0], fd); err != nil {
		fmt.Println("  解析文件描述符失败:", err)
		return
	}
	for _, svc := range fd.GetService() {
		for _, m := range svc.GetMethod() {
			fmt.Printf("    %s (%s) -> (%s)\n", m.GetName(), m.GetInputType(), m.GetOutputType())
		}
	}
}
