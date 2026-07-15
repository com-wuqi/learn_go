package practice

import (
	"errors"
	"time"
)

// ============================================================
// time 包练习（学完 review/time.go 后做）
// ============================================================

// 练习 37：带超时的函数执行
// [TODO] RunWithTimeout 执行 fn()，如果超过 timeout 返回 timeout error
// 提示：goroutine + channel + select + time.After
func RunWithTimeout(fn func() string, timeout time.Duration) (string, error) {
	ch := make(chan string, 1) // 这里需要缓冲
	go func() {
		ch <- fn()
	}()
	select {
	case result := <-ch:
		{
			return result, nil
		}

	case <-time.After(timeout):
		{
			return "", errors.New("timeout")
		}

	}

}

// 练习 38：定时轮询
// [TODO] PollUntil 每隔 interval 检查 fn() 返回的 ok，直到 ok=true 或 timeout
// 返回最终结果和是否成功
// 提示：time.NewTicker + select + time.After
func PollUntil(fn func() (string, bool), interval, timeout time.Duration) (string, bool) {
	ticker := time.NewTicker(interval)
	timeoutCh := time.After(timeout) // 相对于当前时间，如果放进循环则不退出
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if res, ok := fn(); ok {
				return res, ok
			}
		case <-timeoutCh:
			return "", false
		}
	}

}
