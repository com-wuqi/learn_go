package practice

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3"
)

func TestEtcdDiscoverySubscriptionIntegration(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	service := "etcd-subscribe-test"
	prefix := ServicePrefix(service)
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cleanupCancel()
		_, _ = cli.Delete(cleanupCtx, prefix, clientv3.WithPrefix())
	}()

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-1"), "127.0.0.1:9001"); err != nil {
		t.Fatal(err)
	}

	discovery, err := NewServiceDiscovery(ctx, cli, service)
	if err != nil {
		t.Fatal(err)
	}
	defer discovery.Close()

	subscription, err := discovery.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	if subscription == nil {
		t.Fatal("Subscribe() returned nil")
	}
	if got := subscription.Instances["instance-1"]; got != "127.0.0.1:9001" {
		t.Fatalf("initial subscription instance = %q", got)
	}
	if subscription.Revision == 0 {
		t.Fatal("initial subscription revision = 0")
	}

	subscription.Instances["instance-1"] = "mutated"
	if got := discovery.Instances()["instance-1"]; got != "127.0.0.1:9001" {
		t.Fatalf("subscription exposed internal state, got %q", got)
	}

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-2"), "127.0.0.1:9002"); err != nil {
		t.Fatal(err)
	}
	waitForSubscriptionSignal(t, subscription.Updates)
	waitForDiscovery(t, discovery, func(instances map[string]string) bool {
		return instances["instance-2"] == "127.0.0.1:9002"
	})

	subscription.Cancel()
	subscription.Cancel()
	select {
	case _, open := <-subscription.Updates:
		if open {
			t.Fatal("updates channel remained open after Cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("Cancel did not close updates channel")
	}

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-3"), "127.0.0.1:9003"); err != nil {
		t.Fatal(err)
	}
	waitForDiscovery(t, discovery, func(instances map[string]string) bool {
		return instances["instance-3"] == "127.0.0.1:9003"
	})
}

func TestEtcdDiscoverySubscriptionClosesWithDiscovery(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	discovery, err := NewServiceDiscovery(ctx, cli, "etcd-subscribe-close-test")
	if err != nil {
		t.Fatal(err)
	}
	subscription, err := discovery.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	if subscription == nil {
		t.Fatal("Subscribe() returned nil")
	}

	discovery.Close()
	select {
	case _, open := <-subscription.Updates:
		if open {
			t.Fatal("updates channel remained open after discovery.Close")
		}
	case <-time.After(time.Second):
		t.Fatal("discovery.Close did not close subscription updates")
	}
}

func TestEtcdDiscoverySubscribeAfterClose(t *testing.T) {
	d := &ServiceDiscovery{
		instances:   map[string]string{"instance-1": "127.0.0.1:9001"},
		revision:    7,
		subscribers: make(map[chan struct{}]struct{}),
		closed:      true,
	}

	subscription, err := d.Subscribe()
	if !errors.Is(err, ErrDiscoveryClosed) {
		t.Fatalf("Subscribe() error = %v, want ErrDiscoveryClosed", err)
	}
	if subscription != nil {
		t.Fatalf("Subscribe() = %#v on closed discovery, want nil", subscription)
	}
}

func TestEtcdDiscoverySubscribeAfterCloseDoesNotDeadlock(t *testing.T) {
	d := &ServiceDiscovery{
		instances:   make(map[string]string),
		subscribers: make(map[chan struct{}]struct{}),
		closed:      true,
	}
	if subscription, err := d.Subscribe(); !errors.Is(err, ErrDiscoveryClosed) || subscription != nil {
		t.Fatalf("first Subscribe() = (%#v, %v), want ErrDiscoveryClosed", subscription, err)
	}

	returned := make(chan struct{})
	go func() {
		_, _ = d.Subscribe()
		close(returned)
	}()
	select {
	case <-returned:
	case <-time.After(time.Second):
		t.Fatal("Subscribe deadlocked after discovery was already closed")
	}
}

func TestEtcdDiscoveryNotifyConcurrentCancel(t *testing.T) {
	d := &ServiceDiscovery{
		instances:   map[string]string{"instance-1": "127.0.0.1:9001"},
		revision:    1,
		subscribers: make(map[chan struct{}]struct{}),
		runCtx:      context.Background(),
	}
	subscription, err := d.Subscribe()
	if err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			d.notifySubscribers()
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		subscription.Cancel()
	}()
	close(start)
	wg.Wait()
}

func waitForSubscriptionSignal(t *testing.T, updates <-chan struct{}) {
	t.Helper()
	select {
	case _, open := <-updates:
		if !open {
			t.Fatal("subscription updates channel closed unexpectedly")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("subscription did not receive an update")
	}
}
