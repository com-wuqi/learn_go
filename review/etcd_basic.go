package review

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// RunEtcdBasicDemo 演示 etcd Go client v3 基础操作：
// 连接、Put/Get、前缀查询、Watch、Lease（TTL + KeepAlive + Revoke）。
//
// 前置环境（单节点 etcd）：
//
//	docker run -d --name learn-etcd -p 2379:2379 -p 2380:2380 quay.io/coreos/etcd:v3.5.33
func RunEtcdBasicDemo() {
	endpoint := "127.0.0.1:2379"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		fmt.Println("连接 etcd 失败:", err)
		return
	}
	defer cli.Close()
	fmt.Println("✅ 已连接 etcd:", endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// ---------- 1. Put / Get ----------
	fmt.Println("\n--- 1. Put / Get ---")
	putResp, err := cli.Put(ctx, "/demo/hello", "world")
	if err != nil {
		fmt.Println("Put 失败:", err)
		return
	}
	fmt.Printf("Put: /demo/hello = world (revision=%d)\n", putResp.Header.Revision)

	getResp, err := cli.Get(ctx, "/demo/hello")
	if err != nil {
		fmt.Println("Get 失败:", err)
		return
	}
	for _, kv := range getResp.Kvs {
		fmt.Printf("Get: %s = %s (mod_revision=%d)\n", kv.Key, kv.Value, kv.ModRevision)
	}

	// ---------- 2. 前缀查询（服务发现的基石） ----------
	fmt.Println("\n--- 2. 前缀查询 ---")
	// 模拟两个 echo 服务实例：/services/<服务名>/<实例ID> -> 地址
	_, _ = cli.Put(ctx, "/services/echo/1", "127.0.0.1:9001")
	_, _ = cli.Put(ctx, "/services/echo/2", "127.0.0.1:9002")

	listResp, err := cli.Get(ctx, "/services/echo/", clientv3.WithPrefix())
	if err != nil {
		fmt.Println("前缀查询失败:", err)
		return
	}
	for _, kv := range listResp.Kvs {
		fmt.Printf("发现实例: %s -> %s\n", kv.Key, kv.Value)
	}

	// ---------- 3. Watch ----------
	fmt.Println("\n--- 3. Watch（观察变化） ---")
	watchCtx, watchCancel := context.WithCancel(ctx)
	watchCh := cli.Watch(watchCtx, "/services/echo/", clientv3.WithPrefix())
	go func() {
		for resp := range watchCh {
			for _, ev := range resp.Events {
				fmt.Printf("[watch] %s %s = %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	// 等 watch 建立后再写，避免事件丢失
	time.Sleep(300 * time.Millisecond)
	_, _ = cli.Put(ctx, "/services/echo/3", "127.0.0.1:9003")
	time.Sleep(200 * time.Millisecond)
	_, _ = cli.Delete(ctx, "/services/echo/2")
	time.Sleep(200 * time.Millisecond)
	watchCancel()
	fmt.Println("Watch 已取消")

	// ---------- 4. Lease：TTL 租约 ----------
	fmt.Println("\n--- 4. Lease（租约 + KeepAlive） ---")
	// 申请 10 秒租约
	leaseResp, err := cli.Grant(ctx, 10)
	if err != nil {
		fmt.Println("Grant 失败:", err)
		return
	}
	leaseID := leaseResp.ID
	fmt.Printf("申请租约: %x (TTL=10s)\n", leaseID)

	// 绑定租约写入：租约过期 key 自动消失，等价于"实例下线自动摘除"
	_, err = cli.Put(ctx, "/services/echo/tmp", "127.0.0.1:9999", clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Println("带租约写入失败:", err)
		return
	}

	// 开启自动续约，观察 TTL 被"续满"
	keepAliveCh, err := cli.KeepAlive(ctx, leaseID)
	if err != nil {
		fmt.Println("KeepAlive 失败:", err)
		return
	}
	go func() {
		for range keepAliveCh {
			// 每个续约周期都会收到一条响应
		}
	}()

	time.Sleep(2 * time.Second)
	ttlResp, err := cli.TimeToLive(ctx, leaseID)
	if err == nil {
		fmt.Printf("KeepAlive 续约后剩余 TTL: %d 秒\n", ttlResp.TTL)
	}

	// 主动注销（服务优雅下线时调用）
	_, err = cli.Revoke(ctx, leaseID)
	if err != nil {
		fmt.Println("Revoke 失败:", err)
		return
	}
	afterRevoke, err := cli.Get(ctx, "/services/echo/tmp")
	if err == nil && len(afterRevoke.Kvs) == 0 {
		fmt.Println("Revoke 后 key 已自动删除（优雅下线生效）")
	}
}
