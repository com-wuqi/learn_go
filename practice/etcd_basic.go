package practice

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// ============================================================
// 第三阶段 3.3-B：etcd Go client v3 基础
//
// 环境：单节点 etcd 已运行（Docker 容器 learn-etcd），
// endpoint = 127.0.0.1:2379
// 参考演示：review/etcd_basic.go（go run cmd/etcd-demo/main.go）
// ============================================================

// 练习 1：ConnectEtcd
// 用 clientv3.New 连接 etcd，拨号超时设为 3 秒。
// 提示：clientv3.Config{Endpoints: []string{endpoint}, DialTimeout: ...}
// 需要补充 import：time。
func ConnectEtcd(endpoint string) (*clientv3.Client, error) {
	// TODO
	config := clientv3.Config{Endpoints: []string{endpoint}, DialTimeout: 3 * time.Second}
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 练习 2：EtcdPutGet
// 写入 key/value，然后读回该 key 的当前 value。
// 提示：cli.Put 后再 cli.Get；读取 resp.Kvs[0].Value 并转成 string。
func EtcdPutGet(ctx context.Context, cli *clientv3.Client, key, value string) (string, error) {
	// TODO
	_, err := cli.Put(ctx, key, value)
	if err != nil {
		return "", err
	}
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 { // 需要判断
		return "", nil
	}
	val := resp.Kvs[0].Value
	return string(val), nil
}

// 练习 3：EtcdListPrefix
// 列出 prefix 下的全部键值对，返回 map[key]value。
// 这是服务发现的基石：一个服务名一个前缀，实例是前缀下的 key。
// 提示：cli.Get(ctx, prefix, clientv3.WithPrefix())，遍历 resp.Kvs。
func EtcdListPrefix(ctx context.Context, cli *clientv3.Client, prefix string) (map[string]string, error) {
	// TODO
	server := make(map[string]string)
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		server[string(kv.Key)] = string(kv.Value)
	}
	return server, nil
}

// 练习 4：EtcdWatchPrefix
// 监听 prefix 下所有键的变更（PUT / DELETE），把每次变更以
// "PUT /services/echo/1 127.0.0.1:9001" 的格式发到 ch；
// ctx 取消时返回。
// 提示：cli.Watch(ctx, prefix, clientv3.WithPrefix())，
// 遍历 watch channel：for resp := range wch { for _, ev := range resp.Events {...} }，
// ev.Type 是 mvccpb.Event_EventType，用 ev.Type.String()；ev.Kv 是新值。
// 需要补充 import：fmt。
func EtcdWatchPrefix(ctx context.Context, cli *clientv3.Client, prefix string, ch chan<- string) {
	// TODO
	watchCh := cli.Watch(ctx, prefix, clientv3.WithPrefix())
	go func() {
		//for {
		//select {
		//case resp := <-watchCh:
		//	for _, ev := range resp.Events {
		//		ch <- fmt.Sprintf("%s %s %s", ev.Type.String(), ev.Kv.Key, ev.Kv.Value)
		//	} // 这样在 watchCh 关闭后永远拿不到，循环空转,或者监听 ctx.Done()
		for resp := range watchCh {
			for _, ev := range resp.Events {
				ch <- fmt.Sprintf("%s %s %s", ev.Type.String(), ev.Kv.Key, ev.Kv.Value)
			}
		}
		//}
		//}
	}()
}

// 练习 5：EtcdLeasePut
// 申请一个 ttl 秒的租约，把 key=value 绑定到该租约写入，返回租约 ID。
// 租约过期后 key 自动消失——这就是"实例下线自动摘除"的原理。
// 提示：cli.Grant(ctx, ttl) 返回 *clientv3.LeaseGrantResponse，
// 写入时加 clientv3.WithLease(leaseResp.ID)。
func EtcdLeasePut(ctx context.Context, cli *clientv3.Client, key, value string, ttl int64) (clientv3.LeaseID, error) {
	// TODO
	leaseResp, err := cli.Grant(ctx, ttl)
	if err != nil {
		return 0, err
	}
	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		_, _ = cli.Revoke(ctx, leaseResp.ID)
		return 0, err
	}
	return leaseResp.ID, nil
}
