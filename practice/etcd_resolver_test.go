package practice

import (
	"context"
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

func TestEtcdResolverServiceName(t *testing.T) {
	target := resolver.Target{
		URL: url.URL{Scheme: EtcdResolverScheme, Path: "/echo"},
	}
	if got, err := ResolverServiceName(target); err != nil || got != "echo" {
		t.Fatalf("ResolverServiceName() = (%q, %v), want (%q, nil)", got, err, "echo")
	}

	for _, target := range []resolver.Target{
		{URL: url.URL{Scheme: EtcdResolverScheme}},
		{URL: url.URL{Scheme: "other", Path: "/echo"}},
	} {
		if got, err := ResolverServiceName(target); err == nil {
			t.Fatalf("ResolverServiceName(%q) = (%q, nil), want error", target.String(), got)
		}
	}
}

func TestEtcdResolverAddresses(t *testing.T) {
	got := ResolverAddresses(map[string]string{
		"instance-c": "127.0.0.1:9003",
		"instance-a": "127.0.0.1:9001",
		"instance-b": "127.0.0.1:9002",
		"duplicate":  "127.0.0.1:9001",
		"empty":      "",
	})
	want := []resolver.Address{
		{Addr: "127.0.0.1:9001"},
		{Addr: "127.0.0.1:9002"},
		{Addr: "127.0.0.1:9003"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ResolverAddresses() = %#v, want %#v", got, want)
	}
}

func TestEtcdResolverIntegration(t *testing.T) {
	endpoint := os.Getenv("ETCD_ENDPOINT")
	if endpoint == "" {
		t.Skip("set ETCD_ENDPOINT=127.0.0.1:2379 to run etcd integration tests")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	service := "etcd-resolver-test"
	prefix := ServicePrefix(service)
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cleanupCancel()
		_, _ = cli.Delete(cleanupCtx, prefix, clientv3.WithPrefix())
	}()

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-1"), "127.0.0.1:9001"); err != nil {
		t.Fatal(err)
	}

	cc := newRecordingResolverClientConn()
	builder := NewEtcdResolverBuilder(ctx, cli)
	target := resolver.Target{
		URL: url.URL{Scheme: EtcdResolverScheme, Path: "/" + service},
	}
	r, err := builder.Build(target, cc, resolver.BuildOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	select {
	case state := <-cc.states:
		if got, want := resolverAddressStrings(state), []string{"127.0.0.1:9001"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("Build() initial addresses = %v, want %v", got, want)
		}
	default:
		t.Fatal("Build() returned before publishing the initial state")
	}

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-2"), "127.0.0.1:9002"); err != nil {
		t.Fatal(err)
	}
	waitForResolverAddresses(t, cc, []string{"127.0.0.1:9001", "127.0.0.1:9002"})

	if _, err := cli.Delete(ctx, ServiceKey(service, "instance-1")); err != nil {
		t.Fatal(err)
	}
	waitForResolverAddresses(t, cc, []string{"127.0.0.1:9002"})

	r.Close()
	r.Close()
	if _, err := cli.Put(ctx, ServiceKey(service, "instance-3"), "127.0.0.1:9003"); err != nil {
		t.Fatalf("shared etcd client unusable after resolver.Close: %v", err)
	}
}

type recordingResolverClientConn struct {
	resolver.ClientConn
	states chan resolver.State
	errors chan error
}

func newRecordingResolverClientConn() *recordingResolverClientConn {
	return &recordingResolverClientConn{
		states: make(chan resolver.State, 16),
		errors: make(chan error, 16),
	}
}

func (c *recordingResolverClientConn) UpdateState(state resolver.State) error {
	c.states <- state
	return nil
}

func (c *recordingResolverClientConn) ReportError(err error) {
	c.errors <- err
}

func TestEtcdResolverForwardsDiscoveryErrors(t *testing.T) {
	runCtx, runCancel := context.WithCancel(context.Background())
	discoveryDone := make(chan struct{})
	go func() {
		<-runCtx.Done()
		close(discoveryDone)
	}()
	t.Cleanup(func() {
		runCancel()
		<-discoveryDone
	})

	updates := make(chan struct{})
	discovery := &ServiceDiscovery{
		runCtx:    runCtx,
		cancel:    runCancel,
		done:      discoveryDone,
		errors:    make(chan error, 1),
		instances: make(map[string]string),
	}
	subscription := &DiscoverySubscription{
		Updates: updates,
		Cancel:  func() { close(updates) },
	}
	cc := newRecordingResolverClientConn()
	r := &etcdDiscoveryResolver{
		cc:           cc,
		discovery:    discovery,
		subscription: subscription,
		refresh:      make(chan struct{}, 1),
		runDone:      make(chan struct{}),
	}
	go r.Run()

	wantErr := errors.New("watch failed")
	discovery.errors <- wantErr
	select {
	case gotErr := <-cc.errors:
		if !errors.Is(gotErr, wantErr) {
			t.Fatalf("ReportError() = %v, want %v", gotErr, wantErr)
		}
	case <-time.After(time.Second):
		t.Fatal("resolver did not forward discovery error")
	}

	subscription.Cancel()
	select {
	case <-r.runDone:
	case <-time.After(time.Second):
		t.Fatal("resolver run loop did not stop after updates closed")
	}
}

func TestEtcdResolverCloseCancelsOwnedSubscription(t *testing.T) {
	runCtx, runCancel := context.WithCancel(context.Background())
	discoveryDone := make(chan struct{})
	go func() {
		<-runCtx.Done()
		close(discoveryDone)
	}()

	updates := make(chan struct{})
	cancelled := make(chan struct{})
	discovery := &ServiceDiscovery{
		runCtx:      runCtx,
		cancel:      runCancel,
		done:        discoveryDone,
		subscribers: make(map[chan struct{}]struct{}),
		instances:   make(map[string]string),
	}
	subscription := &DiscoverySubscription{
		Updates: updates,
		Cancel: func() {
			close(cancelled)
			close(updates)
		},
	}

	runDone := make(chan struct{})
	close(runDone)
	r := &etcdDiscoveryResolver{
		discovery:    discovery,
		subscription: subscription,
		refresh:      make(chan struct{}, 1),
		runDone:      runDone,
	}
	r.Close()

	select {
	case <-cancelled:
	default:
		t.Fatal("resolver.Close did not cancel its owned subscription")
	}
}

func waitForResolverAddresses(t *testing.T, cc *recordingResolverClientConn, want []string) {
	t.Helper()
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	for {
		select {
		case state := <-cc.states:
			got := resolverAddressStrings(state)
			if reflect.DeepEqual(got, want) {
				return
			}
		case <-timer.C:
			t.Fatalf("resolver addresses did not become %v", want)
		}
	}
}

func resolverAddressStrings(state resolver.State) []string {
	addresses := make([]string, 0, len(state.Addresses))
	for _, address := range state.Addresses {
		addresses = append(addresses, address.Addr)
	}
	return addresses
}
