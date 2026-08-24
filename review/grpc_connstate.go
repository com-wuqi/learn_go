package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// watchConnState 打印连接状态，直到 ctx 结束。
func watchConnState(conn *grpc.ClientConn, ctx context.Context, label string) {
	for {
		state := conn.GetState()
		fmt.Printf("  [%s] state=%s\n", label, state)
		if !conn.WaitForStateChange(ctx, state) {
			return
		}
	}
}

// startGreeterServer 在指定地址起一个 greeter server，返回 server 与实际绑定地址。
func startGreeterServer(addr string) (*grpc.Server, string, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", err
	}
	srv := grpc.NewServer()
	demo.RegisterGreeterServer(srv, &greeterServer{})
	go func() {
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	return srv, ln.Addr().String(), nil
}

// RunConnStateDemo 演示连接状态机：IDLE → CONNECTING → READY →（停服）→ TRANSIENT_FAILURE →（重启）→ READY。
func RunConnStateDemo() {
	fmt.Println("\n=== gRPC 连接状态机演示：停服、TRANSIENT_FAILURE 与恢复 ===")

	srv, addr, err := startGreeterServer("127.0.0.1:0")
	if err != nil {
		fmt.Println("  启动 server 失败:", err)
		return
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	go watchConnState(conn, ctx, "conn")

	fmt.Println("  初始状态（还没 Connect）:", conn.GetState()) // IDLE
	conn.Connect()
	time.Sleep(500 * time.Millisecond)

	client := demo.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		fmt.Println("  第一次 RPC 失败:", err)
		return
	}
	fmt.Println("  第一次 RPC:", resp.GetMessage())

	srv.Stop()
	fmt.Println("  服务已停止；发一个带 3s 超时且 WaitForReady 的 RPC，观察 TRANSIENT_FAILURE...")

	// 等连接完全回到 IDLE，再发一个 WaitForReady 的 RPC：
	// 它会跟着重连尝试走 CONNECTING -> TRANSIENT_FAILURE -> 退避重试，直到超时。
	time.Sleep(300 * time.Millisecond)
	failCtx, failCancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, failErr := demo.NewGreeterClient(conn).SayHello(failCtx, &demo.HelloRequest{Name: "Codex"}, grpc.WaitForReady(true))
	// If waitForReady is false, the RPC will fail immediately
	failCancel()
	fmt.Println("  断连期间的 RPC:", failErr)

	// 同一地址重启，观察客户端自动恢复。
	var newSrv *grpc.Server
	for i := 0; i < 20; i++ {
		newSrv, _, err = startGreeterServer(addr)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		fmt.Println("  重启失败:", err)
		return
	}
	defer newSrv.Stop()
	time.Sleep(1 * time.Second)

	// 恢复 RPC 用 WaitForReady：通道可能还在退避中，fail-fast 会立刻失败，
	// 加上 WaitForReady 让它在重连成功后自动发出。
	resp, err = client.SayHello(ctx, &demo.HelloRequest{Name: "Codex"}, grpc.WaitForReady(true))
	if err != nil {
		fmt.Println("  恢复后 RPC 失败:", err)
		return
	}
	fmt.Println("  恢复后 RPC:", resp.GetMessage())
}
