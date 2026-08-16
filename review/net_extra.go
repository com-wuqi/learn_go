package review

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/netutil"
	"golang.org/x/sync/errgroup"
)

// ============================================================
// 3.1 补充：net.ErrClosed / errgroup / LimitListener
// ============================================================

// ErrClosedDemo 演示如何用 errors.Is(err, net.ErrClosed) 区分
// listener 的正常关闭与真实错误。
func ErrClosedDemo() {
	fmt.Println("\n=== net.ErrClosed：区分正常关闭与真实错误 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	done := make(chan error, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					done <- nil
				} else {
					done <- err
				}
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()

	time.Sleep(50 * time.Millisecond)
	ln.Close()

	if err := <-done; err != nil {
		fmt.Println("  意外错误:", err)
	} else {
		fmt.Println("  正常关闭：Accept 返回 net.ErrClosed，已按 nil 处理")
	}
}

// ErrGroupDemo 演示 errgroup 收集第一个错误，并用 WithContext 取消其余任务。
func ErrGroupDemo() {
	fmt.Println("\n=== errgroup：收集错误 + 协调取消 ===")

	g, ctx := errgroup.WithContext(context.Background())

	// 任务 A：50ms 后失败，触发 ctx 取消
	g.Go(func() error {
		time.Sleep(50 * time.Millisecond)
		return errors.New("task A failed")
	})

	// 任务 B：本应运行 2s，但通过 ctx.Done() 感知取消后退出
	g.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("  task B 感知到取消:", ctx.Err())
			return nil
		case <-time.After(2 * time.Second):
			fmt.Println("  task B 完成")
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("  errgroup 返回第一个错误: %v\n", err)
	}
}

// LimitListenerDemo 演示用 netutil.LimitListener 包装 listener 限制并发连接数。
func LimitListenerDemo() {
	fmt.Println("\n=== LimitListener：限制并发连接数 ===")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}

	limited := netutil.LimitListener(ln, 2)
	defer limited.Close()

	fmt.Println("  用法: netutil.LimitListener(ln, 2)")
	fmt.Println("  超过 2 个并发连接时 Accept 会阻塞，直到有连接释放")
	fmt.Println("  具体并发阻塞行为放在 practice 里验证")
}
