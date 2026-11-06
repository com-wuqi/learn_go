package practice

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

// ============================================================
// 第三阶段 3.3-E2：etcd-backed gRPC resolver
//
// 目标格式：etcd:///<service>，例如 etcd:///echo
// 参考演示：review/grpc_resolver.go（静态 resolver）+
// review/etcd_discover.go（快照 + Watch）
//
// 设计要点：
//   - Builder 只保存 dial 生命周期 ctx 和共享 etcd client，不拥有 client
//   - 每个 resolver 创建并拥有自己的 ServiceDiscovery
//   - ClientConn.UpdateState 必须收到完整地址列表，而非增量
//   - resolver.Close 只关闭自己拥有的订阅和发现器
// ============================================================

const EtcdResolverScheme = "etcd"

type EtcdResolverBuilder struct {
	ctx context.Context
	cli *clientv3.Client
}

func NewEtcdResolverBuilder(ctx context.Context, cli *clientv3.Client) *EtcdResolverBuilder {
	return &EtcdResolverBuilder{ctx: ctx, cli: cli}
}

// ResolverServiceName 从 etcd:///<service> 中解析裸服务名。
//
// 契约：
//   - 返回 target.Endpoint()
//   - target scheme 不是 etcd 或服务名为空时返回 error
//   - 不在这里拼 /services/<service>/，该职责属于 ServiceDiscovery
func ResolverServiceName(target resolver.Target) (string, error) {
	// TODO
	if target.URL.Scheme != EtcdResolverScheme {
		return "", errors.New("unsupported scheme")
	}
	s := target.Endpoint()
	s = strings.TrimSpace(s)
	if s == "/" || s == "" {
		return "", errors.New("invalid target")
	}
	return s, nil
}

// ResolverAddresses 把 ServiceDiscovery 的实例表转换成 gRPC 地址列表。
//
// 契约：
//   - map key 是实例 ID，value 是 host:port；只使用 value
//   - 空地址忽略
//   - 重复地址去重
//   - 结果按 Addr 字典序排序，保证输出稳定
func ResolverAddresses(instances map[string]string) []resolver.Address {
	// TODO
	set := make(map[string]struct{}, len(instances))
	for _, instance := range instances {
		if strings.TrimSpace(instance) != "" {
			set[strings.TrimSpace(instance)] = struct{}{}
		}
	}
	s := make([]string, 0, len(set))
	for instance := range set {
		s = append(s, instance)
	}
	slices.Sort(s)
	resolverAdds := make([]resolver.Address, 0, len(s))
	for _, instance := range s {
		resolverAdds = append(resolverAdds, resolver.Address{Addr: instance})
	}
	return resolverAdds
}

type etcdDiscoveryResolver struct {
	cc           resolver.ClientConn
	discovery    *ServiceDiscovery
	subscription *DiscoverySubscription

	refresh   chan struct{}
	runDone   chan struct{}
	closeOnce sync.Once
}

// Build 创建一个由 ServiceDiscovery 驱动的 resolver。
//
// 契约：
//   - 解析 service 后调用 NewServiceDiscovery，不另写 etcd Get/Watch
//   - Subscribe 失败时关闭已创建的发现器并返回 error
//   - 返回前至少调用一次 cc.UpdateState 发布当前完整地址列表
//   - 后台观察 subscription.Updates 和 discovery.Errors
//   - discovery 的 Updates 关闭时退出后台循环；Close 等待该循环结束
func (b *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// TODO
	server, err := ResolverServiceName(target)
	if err != nil {
		return nil, err
	}
	discovery, err := NewServiceDiscovery(b.ctx, b.cli, server)
	if err != nil {
		return nil, err
	}
	subscription, err := discovery.Subscribe()
	if err != nil {
		discovery.Close()
		return nil, err
	}
	resolve := &etcdDiscoveryResolver{
		cc:           cc,
		discovery:    discovery,
		subscription: subscription,
		refresh:      make(chan struct{}, 1),
		runDone:      make(chan struct{}),
		closeOnce:    sync.Once{},
	}
	resolve.publishState() // 立刻 同步 发布
	go resolve.Run()
	return resolve, nil
}

func (b *EtcdResolverBuilder) Scheme() string {
	return EtcdResolverScheme
}

// ResolveNow 触发一次立即重读；多次调用可Prefix(s, "//")以合并为一次。
func (r *etcdDiscoveryResolver) ResolveNow(resolver.ResolveNowOptions) {
	// TODO
	select {
	case r.refresh <- struct{}{}:
	default:
	}
}

// Close 幂等停止 resolver。
//
// 注意顺序：取消订阅 -> 关闭发现器 -> 等待后台循环退出。
// 不要关闭共享 etcd client。
func (r *etcdDiscoveryResolver) Close() {
	// TODO
	r.closeOnce.Do(func() {
		r.subscription.Cancel()
		// closeSubscribers()：如果手工创建了一个未注册到 discovery.subscribers 的订阅。
		// 此时直接调用 closeSubscribers() 不会关闭这个 Updates
		r.discovery.Close()
		<-r.runDone
		//close(r.refresh) // ResolveNow data race, GC
	})
}

func (r *etcdDiscoveryResolver) Run() {
	defer close(r.runDone)
	updates := r.subscription.Updates
	errorCh := r.discovery.Errors()

	for {
		select {
		case _, ok := <-updates:
			if !ok {
				return
			}
			// publish
			r.publishState()
		case err, ok := <-errorCh:
			if !ok {
				errorCh = nil
				continue
			}
			if err != nil {
				r.cc.ReportError(err)
			}
		case <-r.refresh:
			// publish
			r.publishState()
		}
	}
}

func (r *etcdDiscoveryResolver) publishState() {
	state := resolver.State{Addresses: ResolverAddresses(r.discovery.Instances())}
	err := r.cc.UpdateState(state)
	if err != nil {
		r.cc.ReportError(err)
		return
	}
}

// ------------------------------------------------------------
// 自检：
//   - etcd:///echo 解析为 echo；空 target/错误 scheme 明确报错
//   - 初始快照立即发布，地址排序且去重
//   - etcd PUT/DELETE 后 ClientConn 收到完整新列表
//   - discovery.Errors 转发到 ClientConn.ReportError
//   - ResolveNow 并发调用不阻塞、不泄漏 goroutine
//   - Close 幂等，等待后台退出，且仍可继续使用共享 etcd client
// ------------------------------------------------------------
