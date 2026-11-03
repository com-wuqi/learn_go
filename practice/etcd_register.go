package practice

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// ============================================================
// 第三阶段 3.3-C：服务注册（etcd）
//
// 前置：3.3-B（practice/etcd_basic.go）已完成
// 环境：etcd 容器 learn-etcd 运行中，endpoint = 127.0.0.1:2379
// 参考演示：review/etcd_basic.go 的 Lease + KeepAlive + Revoke 部分
//
// 设计要点：
//   - key 格式 /services/<服务名>/<实例ID>，实例 ID 每个进程唯一
//   - 注册 = Grant(ttl) + Put(..., WithLease(id))，进程猝死后 key 随租约自动消失
//   - 心跳 = KeepAlive(ctx, id)，ctx 必须是进程生命周期级别（不能带短超时）
//   - 优雅注销 = Revoke(id)；cli.Close() 不会撤销租约
// ============================================================

// 练习 1：ServiceKey
// 生成实例注册 key，格式：/services/<service>/<instanceID>
// 提示：fmt.Sprintf 即可；约定 service 与 instanceID 均非空。
func ServiceKey(service, instanceID string) string {
	// TODO
	return fmt.Sprintf("/services/%s/%s", service, instanceID)
}

// 练习 2：RegisterService
// 注册一个服务实例：
//  1. 申请 ttl 秒租约（cli.Grant）
//  2. 把 key=ServiceKey(service, instanceID)、value=addr 写入并绑定该租约（clientv3.WithLease）
//
// 返回租约 ID，调用方用它续约（练习 3）或注销（练习 4）。
// 提示：写入失败时记得 Revoke 掉已申请的租约，避免留下没人管的租约（3.3-B 练习 5 做过）。
func RegisterService(ctx context.Context, cli *clientv3.Client, service, instanceID, addr string, ttl int64) (clientv3.LeaseID, error) {
	// TODO
	leaseResp, err := cli.Grant(ctx, ttl)
	if err != nil {
		return 0, err
	}
	key := ServiceKey(service, instanceID)
	_, err = cli.Put(ctx, key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		_, err2 := cli.Revoke(ctx, leaseResp.ID)
		if err2 != nil {
			return 0, fmt.Errorf("failed: %w %w", err, err2)
		} // %w better
		return 0, err
	}
	return leaseResp.ID, nil
}

// 练习 3：KeepServiceAlive
// 在后台持续给 leaseID 续约，直到 ctx 取消。
// 返回的 channel 在续约循环结束后关闭，语义约定：
//   - ctx 取消 → 正常下线，不发错误
//   - ctx 未取消但续约通道结束（租约被撤销/过期、etcd 不可达）→ 发送一个 error 后关闭
//
// 提示：cli.KeepAlive(ctx, leaseID) 返回 <-chan *clientv3.LeaseKeepAliveResponse，
// range 它即可（ctx 取消会关闭该 channel）；循环结束后用 ctx.Err() 判断是正常下线还是异常丢失。
func KeepServiceAlive(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) <-chan error {
	// TODO
	resp, err := cli.KeepAlive(ctx, leaseID)
	errCh := make(chan error, 1)
	if err != nil {
		errCh <- err
		close(errCh)
		return errCh
	}
	go func() {
		defer close(errCh)
		for range resp {
		}
		if ctx.Err() == nil {
			errCh <- fmt.Errorf("keepalive 结束但 ctx 未取消，租约 %x 可能已丢失", leaseID)
		}
	}()

	return errCh
}

// 练习 4：DeregisterService
// 优雅注销：撤销租约，绑定在它上面的 key 会一并删除。
// 提示：cli.Revoke(ctx, leaseID)。
func DeregisterService(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) error {
	// TODO
	_, err := cli.Revoke(ctx, leaseID)
	return err
}

// 练习 5（进阶）：RegisterServiceWithHeartbeat
// 组合练习 2/3/4，把"注册 + 后台续约 + 优雅注销"打包成一个可用的注册器：
//   - 返回 leaseID 和一个 stop 函数
//   - stop 必须幂等：第一次调用停止续约并 Revoke（返回其错误），之后调用返回 nil
//
// 提示：用 context.WithCancel 控制续约生命周期，用 sync.Once（或 CAS 标记）保证只执行一次；
// 别忘了 stop 之后要把 KeepServiceAlive 返回的 error channel 处理掉（正常取消不算错误）。
func RegisterServiceWithHeartbeat(ctx context.Context, cli *clientv3.Client, service, instanceID, addr string, ttl int64) (clientv3.LeaseID, func() error, error) {
	// TODO
	key := fmt.Sprintf("/services/%s/%s", service, instanceID)
	ctx, cancel := context.WithCancel(ctx)
	lease, err := cli.Grant(ctx, ttl)
	if err != nil {
		cancel()
		return 0, nil, err
	}
	resp, err := cli.KeepAlive(ctx, lease.ID)
	if err != nil {
		cancel()
		return 0, nil, err
	}
	_, err = cli.Put(ctx, key, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		cancel()
		ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel3()
		_, err2 := cli.Revoke(ctx3, lease.ID) // cleanUp 使用独立ctx
		if err2 != nil {
			return 0, nil, errors.Join(err, err2)
		}
		return 0, nil, err
	}
	errCh := make(chan error, 1)
	var once sync.Once
	go func() {
		defer close(errCh)
		for range resp {
		}
		if ctx.Err() == nil {
			errCh <- fmt.Errorf("keepalive 结束但 ctx 未取消，租约 %x 可能已丢失", lease.ID)
		}
	}()
	stop := func() error {
		var err2 error
		var err error
		once.Do(func() { // 保证幂等
			cancel()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			_, err = cli.Revoke(ctx2, lease.ID)
			err2, _ = <-errCh
		})
		return errors.Join(err, err2)
	}
	return lease.ID, stop, nil
}
