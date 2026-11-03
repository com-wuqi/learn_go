package practice

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

// ============================================================
// 第三阶段 3.3-D：服务发现（本地实例列表 + Watch 增量维护）
//
// 前置：3.3-B（etcd client）与 3.3-C（服务注册）已完成
// 环境：etcd 容器 learn-etcd 运行中，endpoint = 127.0.0.1:2379
//
// 总体设计：
//  1. 先用一次 linearizable Get 拿全量快照，并记录 Header.Revision
//  2. 再从 snapshotRevision+1 开始 Watch，重放快照之后发生的 PUT/DELETE
//  3. 每次事件只维护本地 map；读取者通过复制值拿到稳定快照
//  4. Watch 被压缩等原因取消时，重新 Get 全量快照，再原子替换本地状态
//
// 关键一致性点：不能让 Watch 从 revision=0 开始后才做 Get，
// 否则两者之间的变更可能漏掉。应始终“先 Get，再 Watch(rev+1)”。
// ============================================================

// 练习 1：ServicePrefix
// 返回某个服务全部实例的 etcd 前缀，格式：/services/<service>/
// 注意结尾必须有 "/"，否则 /services/echo 前缀也可能匹配到别的服务。
func ServicePrefix(service string) string {
	// TODO
	return fmt.Sprintf("/services/%s/", service)
}

// 练习 2：LoadServiceSnapshot
// 一次性读取 service 的全部在线实例，并返回快照对应的 etcd revision。
//
// 返回值 map 的 key 是实例 ID，value 是实例地址，例如：
//
//	"instance-1" -> "127.0.0.1:9001"
//
// 提示：
//   - cli.Get(ctx, prefix, clientv3.WithPrefix())
//   - Get 默认是 linearizable 的；resp.Header.Revision 是后续 Watch 的起点
//   - key 去掉 ServicePrefix 后得到实例 ID；如果 key 不合法，可以选择忽略
//   - 返回空 map 而不是 nil，调用方维护空集合会更方便
func LoadServiceSnapshot(ctx context.Context, cli *clientv3.Client, service string) (map[string]string, int64, error) {
	// TODO
	key := ServicePrefix(service)
	aMap := make(map[string]string)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return aMap, 0, err
	}
	for _, kv := range resp.Kvs {
		if !strings.HasPrefix(string(kv.Key), key) || len(kv.Key) < len(key) {
			continue
		}
		sid, _ := strings.CutPrefix(string(kv.Key), key)
		if len(sid) == 0 || strings.Contains(sid, "/") {
			continue
		}
		aMap[sid] = string(kv.Value)
	}
	return aMap, resp.Header.Revision, nil
}

// 练习 3：ApplyServiceEvent
// 把一条 Watch 事件应用到本地实例表。
//   - PUT：写入或覆盖 instances[instanceID] = addr
//   - DELETE：删除 instances[instanceID]
//
// 提示：
//   - 事件类型是 mvccpb.PUT / mvccpb.DELETE，可通过 ev.Type 判断
//   - ev.Kv 理论上非空，但生产代码应防御 nil
//   - 应验证 key 确实位于 ServicePrefix(service) 下
//   - map 的并发保护由调用方负责；本函数只处理一条事件
func ApplyServiceEvent(instances map[string]string, service string, ev *clientv3.Event) {
	// TODO
	if ev.Kv == nil {
		return
	}
	insId, ok := strings.CutPrefix(string(ev.Kv.Key), ServicePrefix(service))
	if !ok || len(insId) == 0 || strings.Contains(insId, "/") {
		return
	}
	switch ev.Type {
	case mvccpb.Event_PUT:
		{
			instances[insId] = string(ev.Kv.Value)
		}
	case mvccpb.Event_DELETE:
		{
			delete(instances, insId)
		}
	default:
		return
	}
}

// ServiceDiscovery 维护一个服务的本地实例视图。
//
// 字段只给出一组可行设计，可以在实现时调整；若调整，请保持
// Instances/Revision/Errors/Close 的对外契约不变。
type ServiceDiscovery struct {
	cli     *clientv3.Client
	service string

	mu        sync.RWMutex
	instances map[string]string
	revision  int64 // 对外 revision

	subMu       sync.Mutex
	subscribers map[chan struct{}]struct{}
	closed      bool

	runCtx    context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
	done      chan struct{}
	errors    chan error
}

// DiscoverySubscription 是服务发现状态的一份初始视图和后续变更通知。
//
// 契约：
//   - Instances 是独立副本，注册订阅时已经读取
//   - Updates 使用容量为 1 的 channel，同一批尚未消费的变更可以合并
//   - 收到通知表示“状态可能变化”，消费者应重新调用 discovery.Instances()
//   - Cancel 幂等，只取消这个订阅，不停止 ServiceDiscovery
//   - ServiceDiscovery 退出时必须关闭所有订阅的 Updates
type DiscoverySubscription struct {
	Instances map[string]string
	Revision  int64
	Updates   <-chan struct{}
	Cancel    func()
}

// 练习 4：NewServiceDiscovery
// 创建并启动服务发现器：
//  1. 调用 LoadServiceSnapshot 得到初始实例表和 revision
//  2. 从 snapshotRevision+1 开始后台 Watch
//  3. 返回的 discovery 可被多个 goroutine 并发读取
//
// 契约：
//   - 初始化快照失败时返回 error，不启动后台 goroutine
//   - ctx 是 discovery 的生命周期；ctx 取消后 Watch 与 done 都要结束
//   - cli 由调用方拥有，ServiceDiscovery 不负责 cli.Close()
func NewServiceDiscovery(ctx context.Context, cli *clientv3.Client, service string) (*ServiceDiscovery, error) {
	// TODO
	ctx2, cancel2 := context.WithCancel(ctx)
	insMap, revision, err := LoadServiceSnapshot(ctx2, cli, service)
	if err != nil {
		cancel2()
		return nil, err
	}
	sDiscovery := &ServiceDiscovery{
		cli:         cli,
		service:     service,
		mu:          sync.RWMutex{},
		instances:   insMap,
		revision:    revision,
		runCtx:      ctx2,
		cancel:      cancel2,
		closeOnce:   sync.Once{},
		done:        make(chan struct{}),
		errors:      make(chan error, 1),
		subscribers: make(map[chan struct{}]struct{}),
	}
	go sDiscovery.run(revision + 1)
	return sDiscovery, nil
}

// 练习 5：后台维护本地状态
// 读取 watchCtx 上的 WatchResponse：
//   - 正常响应：逐条 ApplyServiceEvent；更新 revision；必要时通知等待者
//   - resp.Canceled：记录 resp.Err()，退出当前 Watch
//   - ErrCompacted：必须重新 LoadServiceSnapshot，再根据新 revision 重连 Watch
//   - 其他不可恢复错误：可以在 Errors 上报告后重建 Watch；注意退避，避免忙循环
//   - 主 ctx 取消：正常退出，关闭 errors 与 done，不报告 cancellation 为故障
//
// 普通网络断连通常由 etcd client 自己重连；不要另起 goroutine 轮询 etcd。
// 唯一必须显式处理的常见一致性故障是 Watch 起点被压缩（ErrCompacted）。
//
// 并发约束：
//   - instances/revision 只能在该后台 goroutine 或持有写锁时修改
//   - 重新同步时先构建新 map，再一次性加锁替换，避免读取者看到半成品
//   - errors 建议 buffered + 非阻塞发送，防止无人消费时卡死发现器
func (d *ServiceDiscovery) run(startRevision int64) {
	// TODO
	defer close(d.done)        // 通知 Close()
	defer close(d.errors)      // 错误
	defer d.closeSubscribers() // 关闭 Subscribers
	revision := startRevision  // 下一次使用的
	needToResync := false
	for {
		if needToResync {
			insMap, newRevision, err := LoadServiceSnapshot(d.runCtx, d.cli, d.service)
			if err != nil {
				if d.runCtx.Err() != nil {
					return // 捕获退出
				}
				d.reportError(fmt.Errorf("reload snapshot for service %q: %w", d.service, err))
				select {
				case <-time.After(200 * time.Millisecond):
					{
						continue
					}
				case <-d.runCtx.Done():
					return
				}
			}
			d.mu.Lock()
			d.instances = insMap
			d.revision = newRevision
			d.mu.Unlock()
			d.notifySubscribers() // 重建后通知
			newRevision += 1      // 锁外固定revision为接下来watch的
			revision = newRevision
			needToResync = false
			continue
		}
		watch := d.cli.Watch(d.runCtx, ServicePrefix(d.service), clientv3.WithPrefix(), clientv3.WithRev(revision))
		for resp := range watch {
			if d.runCtx.Err() != nil {
				return // 捕获退出
			}
			if err := resp.Err(); err != nil {
				d.reportError(fmt.Errorf("watch service %q from revision %d: %w",
					d.service, revision, err))
				needToResync = true
				break
			}
			for _, ev := range resp.Events {
				d.mu.Lock()
				ApplyServiceEvent(d.instances, d.service, ev)
				d.mu.Unlock()
			}
			if len(resp.Events) > 0 {
				d.notifySubscribers() // 存在事件通知
			}
			if resp.Header != nil {
				d.mu.Lock()
				d.revision = resp.Header.GetRevision()
				d.mu.Unlock()
				revision = resp.Header.GetRevision() + 1
			}
		}
		if d.runCtx.Err() != nil {
			return
		}
	}
}

func (d *ServiceDiscovery) reportError(err error) {
	if err == nil {
		return
	}
	select {
	case d.errors <- err:
	default:
	}
}

// ErrDiscoveryClosed Discovery 被关闭后尝试注册服务发现状态变更
var ErrDiscoveryClosed = errors.New("service discovery is closed")

// Subscribe 注册服务发现状态变更。
// 提示：先注册 subscriber，再读取初始快照，可避免两次操作之间漏掉通知；
// 这段时间内收到的一次额外通知是无害的。不要同时持有 d.mu 和 d.subMu。
func (d *ServiceDiscovery) Subscribe() (*DiscoverySubscription, error) {
	// TODO
	// 注册
	d.subMu.Lock()
	if d.closed || d.runCtx.Err() != nil {
		// closed 或者 runCtx 已经结束
		d.subMu.Unlock()
		return nil, ErrDiscoveryClosed
	}
	updates := make(chan struct{}, 1)
	d.subscribers[updates] = struct{}{}
	d.subMu.Unlock()

	// 数据
	d.mu.RLock()
	instances := make(map[string]string)
	maps.Copy(instances, d.instances)
	revision := d.revision
	d.mu.RUnlock()

	cancel := func() {
		d.subMu.Lock()
		defer d.subMu.Unlock()
		if _, ok := d.subscribers[updates]; !ok {
			return
		}
		delete(d.subscribers, updates)
		close(updates)
	}
	return &DiscoverySubscription{
		Instances: instances,
		Revision:  revision,
		Updates:   updates,
		Cancel:    cancel,
	}, nil
}

// notifySubscribers 向所有订阅者发送“状态可能变化”信号。
// 提示：每个 channel 都使用非阻塞发送；同一订阅已经有未消费信号时直接合并。
func (d *ServiceDiscovery) notifySubscribers() {
	// TODO
	// 需要加锁，回收会删除条目
	d.subMu.Lock()
	defer d.subMu.Unlock()
	for update := range d.subscribers {
		select {
		case update <- struct{}{}:
		default:
		}
	}
}

// closeSubscribers 在发现器退出时关闭并清空所有订阅。
// 提示：关闭 channel 和标记 closed 必须在同一把锁内完成。
func (d *ServiceDiscovery) closeSubscribers() {
	// TODO
	d.subMu.Lock()
	defer d.subMu.Unlock()
	d.closed = true // 关闭标记
	for update := range d.subscribers {
		close(update)
		delete(d.subscribers, update)
	}
}

// Instances 返回当前实例表的副本，调用方可以安全修改返回值。
// 提示：必须在读锁内复制 map；不能直接返回 d.instances。
func (d *ServiceDiscovery) Instances() map[string]string {
	// TODO
	insMap := make(map[string]string)
	d.mu.RLock()
	maps.Copy(insMap, d.instances)
	d.mu.RUnlock()
	return insMap
}

// Revision 返回本地状态已经应用到的 etcd revision。
func (d *ServiceDiscovery) Revision() int64 {
	// TODO
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.revision
}

// Errors 返回后台发现过程中需要调用方知晓的错误。
// channel 应在后台 goroutine 退出时关闭。
func (d *ServiceDiscovery) Errors() <-chan error {
	// TODO
	return d.errors
}

// Close 停止后台 Watch 并等待 goroutine 退出；必须幂等。
// 不要关闭 cli，因为同一个 client 可能还被注册器或其他组件使用。
func (d *ServiceDiscovery) Close() {
	// TODO
	d.closeOnce.Do(func() {
		d.cancel()
	})
	<-d.done // 等待
}

// ------------------------------------------------------------
// 自检清单（实现后逐项验证）：
//   - 启动时已有实例出现在 Instances 中
//   - 新实例 PUT 后被加入，地址更新后覆盖，DELETE 后移除
//   - 快照完成与 Watch 建立之间的并发写入不会漏掉
//   - 不存在的服务返回空表而不是 error
//   - Instances 返回副本，外部修改不会污染内部状态
//   - 快照失败时 NewServiceDiscovery 不泄漏 goroutine
//   - context 取消/Close 后后台 goroutine 正常退出，Close 可重复调用
//   - Watch 因 compression 取消后能通过重同步恢复到最新状态
// ------------------------------------------------------------
// ETCD_ENDPOINT=127.0.0.1:2379 go test ./practice -run EtcdDiscoverySubscription -v
