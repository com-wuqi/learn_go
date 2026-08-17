package practice

import (
	"context"
	"errors"
	"io"
	"net"

	"golang.org/x/net/netutil"
	"golang.org/x/sync/errgroup"
)

// ============================================================
// 3.1 补充练习：net.ErrClosed / errgroup / LimitListener
// ============================================================

// 练习 net.ErrClosed：AcceptLoop
// 循环 Accept 并处理每个连接。
// 当 listener 被正常关闭时（errors.Is(err, net.ErrClosed)），返回 nil；
// 其他错误原样返回。
func AcceptLoop(ln net.Listener) error {
	// TODO
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go func(conn net.Conn) {
			defer conn.Close()
			io.Copy(conn, conn)
		}(conn)
	}
}

// 练习 errgroup：并发执行多个任务，任一失败即取消其余任务并返回第一个错误。
// 每个 task 会收到 errgroup 派生出的 ctx，需要自行监听 ctx.Done() 实现协作取消。
// 提示: errgroup.WithContext
func RunWithErrGroup(ctx context.Context, tasks []func(ctx context.Context) error) error {
	// TODO
	group, errCtx := errgroup.WithContext(ctx)
	for _, task := range tasks {
		group.Go(func() error {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			default:
				return task(errCtx)
			}
		})
	}
	return group.Wait()
}

// 练习 LimitListener：包装 listener，限制最大并发连接数。
// 提示: netutil.LimitListener
func NewLimitedListener(ln net.Listener, maxConns int) net.Listener {
	// TODO
	limiter := netutil.LimitListener(ln, maxConns)
	return limiter
}
