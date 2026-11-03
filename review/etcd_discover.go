package review

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/client/v3"
)

// reviewServiceDiscovery 是 3.3-D 的完整演示实现。
// 它对应 practice/etcd_discover.go 的五个练习：
//  1. prefix 统一使用 /services/<service>/
//  2. 初始状态来自一次 linearizable Get 快照
//  3. PUT/DELETE 增量维护 instanceID -> addr
//  4. 从 snapshotRevision+1 启动 Watch，读取端拿 map 副本
//  5. Watch 失效时重同步；ctx 取消或 Close 后退出
type reviewServiceDiscovery struct {
	cli    *clientv3.Client
	prefix string

	mu        sync.RWMutex
	instances map[string]string
	revision  int64

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func newReviewServiceDiscovery(parent context.Context, cli *clientv3.Client, service string) (*reviewServiceDiscovery, error) {
	prefix := "/services/" + service + "/"
	instances, revision, err := reviewServiceSnapshot(parent, cli, prefix)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(parent)
	d := &reviewServiceDiscovery{
		cli:       cli,
		prefix:    prefix,
		instances: instances,
		revision:  revision, // Header.Revision 仍是集群全局 revision。所以 revision 可能出现跳跃，不代表监听了其他前缀的事件
		ctx:       ctx,
		cancel:    cancel,
		done:      make(chan struct{}),
	}
	go d.run(revision)
	return d, nil
}

func reviewServiceSnapshot(ctx context.Context, cli *clientv3.Client, prefix string) (map[string]string, int64, error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, 0, err
	}

	instances := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		instanceID := string(kv.Key[len(prefix):])
		if instanceID != "" {
			instances[instanceID] = string(kv.Value)
		}
	}
	return instances, resp.Header.Revision, nil
}

func (d *reviewServiceDiscovery) run(startRevision int64) {
	defer close(d.done)
	revision := startRevision

	for {
		watchCh := d.cli.Watch(
			d.ctx,
			d.prefix,
			clientv3.WithPrefix(),
			clientv3.WithRev(revision+1), // 保存的revision+1
		)
		resync := false

		for resp := range watchCh {
			if resp.Canceled {
				if err := resp.Err(); err != nil && !errors.Is(err, context.Canceled) {
					fmt.Printf("[discovery] Watch 中断: %v\n", err)
				}
				resync = true
				break
			}
			d.apply(resp.Events)
			if resp.Header != nil {
				revision = resp.Header.Revision
			}
		}

		if d.ctx.Err() != nil {
			return
		}
		if !resync {
			fmt.Println("[discovery] Watch channel 意外关闭，准备重同步")
		}

		instances, newRevision, err := reviewServiceSnapshot(d.ctx, d.cli, d.prefix)
		if err != nil {
			if d.ctx.Err() != nil {
				return
			}
			fmt.Printf("[discovery] 快照重同步失败，稍后重试: %v\n", err)
			select {
			case <-time.After(250 * time.Millisecond):
			case <-d.ctx.Done():
				return
			}
			continue
		}

		d.mu.Lock() // atom update
		d.instances = instances
		d.revision = newRevision
		d.mu.Unlock()
		revision = newRevision
		fmt.Printf("[discovery] 已重同步到 revision=%d\n", revision)
	}
}

func (d *reviewServiceDiscovery) apply(events []*clientv3.Event) {
	if len(events) == 0 {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, ev := range events {
		if ev.Kv == nil {
			continue
		}

		key := string(ev.Kv.Key)
		if len(key) <= len(d.prefix) || key[:len(d.prefix)] != d.prefix {
			// 客户端防御性校验
			continue
		}
		// 也可以通过尝试切割前缀来校验
		eventInstanceId, ok := strings.CutPrefix(key, d.prefix)
		if !ok || eventInstanceId == "" ||
			strings.Contains(eventInstanceId, "\\") || strings.Contains(eventInstanceId, "/") {
			continue
		}

		instanceID := key[len(d.prefix):]
		switch ev.Type {
		case clientv3.EventTypePut:
			d.instances[instanceID] = string(ev.Kv.Value)
		case clientv3.EventTypeDelete:
			delete(d.instances, instanceID)
		}
	}
}

func (d *reviewServiceDiscovery) snapshot() (map[string]string, int64) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	instances := make(map[string]string, len(d.instances))
	maps.Copy(instances, d.instances)
	return instances, d.revision
}

func (d *reviewServiceDiscovery) close() {
	d.cancel()
	<-d.done
}

// RunEtcdDiscoverDemo 演示服务发现的完整状态流：
// 全量快照 -> Watch PUT/DELETE -> 租约过期自动摘除 -> 压缩后重同步。
//
// 注意：为了真实触发 ErrCompacted，本演示会压缩 etcd 历史 revision。
// 只应在 learn-etcd 学习容器中运行，不要指向生产环境。
//
// 运行：
//
//	go run cmd/etcd-demo/main.go
func RunEtcdDiscoverDemo() {
	const endpoint = "127.0.0.1:2379"
	service := fmt.Sprintf("review-echo-%d", time.Now().UnixNano())
	prefix := "/services/" + service + "/"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		fmt.Println("加载 review/etcd_discover.go 失败:", err)
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cleanupCancel()
		resp, err := cli.Delete(cleanupCtx, prefix, clientv3.WithPrefix())
		if err != nil {
			fmt.Printf("[cleanup] 删除演示实例失败: %v\n", err)
			return
		}
		fmt.Printf("[cleanup] 已删除 %d 个演示实例\n", resp.Deleted)
	}()

	if _, err := cli.Put(ctx, prefix+"instance-1", "127.0.0.1:9101"); err != nil {
		fmt.Println("准备实例 1 失败:", err)
		return
	}
	if _, err := cli.Put(ctx, prefix+"instance-2", "127.0.0.1:9102"); err != nil {
		fmt.Println("准备实例 2 失败:", err)
		return
	}

	// 对应练习 1-4：先拿快照，再从 revision+1 监听。
	discovery, err := newReviewServiceDiscovery(ctx, cli, service)
	if err != nil {
		fmt.Println("启动服务发现失败:", err)
		return
	}
	defer discovery.close()

	instances, revision := discovery.snapshot()
	fmt.Printf("\n--- 3.3-D 服务发现（初始快照 revision=%d）---\n", revision)
	printReviewInstances(instances)

	// 对应练习 3 和 5：Watch 增量维护，不重新扫描整个前缀。
	if _, err := cli.Put(ctx, prefix+"instance-3", "127.0.0.1:9103"); err != nil {
		fmt.Println("新增实例失败:", err)
		return
	}
	wantAfterAdd := map[string]string{
		"instance-1": "127.0.0.1:9101",
		"instance-2": "127.0.0.1:9102",
		"instance-3": "127.0.0.1:9103",
	}
	if !waitForReviewInstances(ctx, discovery, wantAfterAdd) {
		fmt.Println("等待新增实例事件超时")
		return
	}
	fmt.Println("PUT 新实例 → 本地视图已增量加入")
	printReviewInstances(wantAfterAdd)

	if _, err := cli.Put(ctx, prefix+"instance-1", "127.0.0.1:9191"); err != nil {
		fmt.Println("更新实例地址失败:", err)
		return
	}
	wantAfterUpdate := map[string]string{
		"instance-1": "127.0.0.1:9191",
		"instance-2": "127.0.0.1:9102",
		"instance-3": "127.0.0.1:9103",
	}
	if !waitForReviewInstances(ctx, discovery, wantAfterUpdate) {
		fmt.Println("等待地址更新事件超时")
		return
	}
	fmt.Println("PUT 同实例新地址 → 本地视图已覆盖")
	printReviewInstances(wantAfterUpdate)

	if _, err := cli.Delete(ctx, prefix+"instance-2"); err != nil {
		fmt.Println("删除实例失败:", err)
		return
	}
	wantAfterDelete := map[string]string{
		"instance-1": "127.0.0.1:9191",
		"instance-3": "127.0.0.1:9103",
	}
	if !waitForReviewInstances(ctx, discovery, wantAfterDelete) {
		fmt.Println("等待删除实例事件超时")
		return
	}
	fmt.Println("DELETE 实例 → 本地视图已增量移除")
	printReviewInstances(wantAfterDelete)

	demonstrateReviewCompactionRecovery(ctx, cli, prefix)
	fmt.Println("\n3.3-D 服务发现演示结束")
}

func demonstrateReviewCompactionRecovery(ctx context.Context, cli *clientv3.Client, prefix string) {
	fmt.Println("\n--- 故障容错：Watch 起点被压缩 ---")

	before, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("读取压缩前 revision 失败:", err)
		return
	}
	staleRevision := before.Header.Revision

	putResp, err := cli.Put(ctx, prefix+"compaction-probe", "127.0.0.1:9199")
	if err != nil {
		fmt.Println("写入压缩测试实例失败:", err)
		return
	}
	if _, err := cli.Compact(ctx, putResp.Header.Revision); err != nil {
		fmt.Println("压缩 etcd historical revisions 失败:", err)
		return
	}

	// 压缩后故意从一个旧 revision 开始 Watch；服务端会返回 ErrCompacted。
	watchCtx, watchCancel := context.WithTimeout(ctx, 3*time.Second)
	defer watchCancel()
	watchCh := cli.Watch(
		watchCtx,
		prefix,
		clientv3.WithPrefix(),
		clientv3.WithRev(staleRevision),
	)

	for {
		resp, ok := <-watchCh
		if !ok {
			fmt.Println("旧 revision Watch 已关闭，但未收到错误响应")
			return
		}
		if !resp.Canceled {
			continue
		}
		if !errors.Is(resp.Err(), rpctypes.ErrCompacted) {
			fmt.Printf("旧 revision Watch 被其他错误取消: %v\n", resp.Err())
			return
		}
		fmt.Printf("旧 revision Watch 被取消: %v\n", resp.Err())
		break
	}

	recovered, revision, err := reviewServiceSnapshot(ctx, cli, prefix)
	if err != nil {
		fmt.Println("压缩后重新获取快照失败:", err)
		return
	}
	fmt.Printf("重新获取全量快照 revision=%d → 下次 Watch 将从 revision=%d 开始\n", revision, revision+1)
	printReviewInstances(recovered)
}

func waitForReviewInstances(ctx context.Context, d *reviewServiceDiscovery, want map[string]string) bool {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		got, _ := d.snapshot()
		if reflect.DeepEqual(got, want) {
			return true
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return false
		}
	}
}

func printReviewInstances(instances map[string]string) {
	ids := make([]string, 0, len(instances))
	for id := range instances {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	if len(ids) == 0 {
		fmt.Println("  (无在线实例)")
		return
	}
	for _, id := range ids {
		fmt.Printf("  %s -> %s\n", id, instances[id])
	}
}
