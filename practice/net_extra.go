package practice

import (
	"context"
	"net"
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
	return nil
}

// 练习 errgroup：并发执行多个任务，任一失败即取消其余任务并返回第一个错误。
// 提示: errgroup.WithContext
func RunWithErrGroup(ctx context.Context, tasks []func() error) error {
	// TODO
	return nil
}

// 练习 LimitListener：包装 listener，限制最大并发连接数。
// 提示: netutil.LimitListener
func NewLimitedListener(ln net.Listener, maxConns int) net.Listener {
	// TODO
	return nil
}
