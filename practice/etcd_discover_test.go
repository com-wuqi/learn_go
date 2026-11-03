package practice

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

func TestEtcdDiscoveryServicePrefixAndEvents(t *testing.T) {
	if got, want := ServicePrefix("echo"), "/services/echo/"; got != want {
		t.Fatalf("ServicePrefix() = %q, want %q", got, want)
	}

	instances := make(map[string]string)
	put := &clientv3.Event{
		Type: mvccpb.Event_PUT,
		Kv: &mvccpb.KeyValue{
			Key:   []byte(ServiceKey("echo", "instance-1")),
			Value: []byte("127.0.0.1:9001"),
		},
	}
	ApplyServiceEvent(instances, "echo", put)
	if got := instances["instance-1"]; got != "127.0.0.1:9001" {
		t.Fatalf("PUT result = %q", got)
	}

	ApplyServiceEvent(instances, "other", put)
	if _, exists := instances["instance-1"]; !exists {
		t.Fatal("non-matching event removed an instance")
	}

	ApplyServiceEvent(instances, "echo", &clientv3.Event{
		Type: mvccpb.Event_DELETE,
		Kv:   put.Kv,
	})
	if _, exists := instances["instance-1"]; exists {
		t.Fatal("DELETE did not remove the instance")
	}
}

func TestEtcdDiscoveryErrorReportingIsNonBlocking(t *testing.T) {
	d := &ServiceDiscovery{errors: make(chan error, 1)}
	d.reportError(errors.New("first"))

	done := make(chan struct{})
	go func() {
		d.reportError(errors.New("second"))
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("reportError blocked after its buffer filled")
	}
}

func TestEtcdDiscoveryCloseWaitsForDone(t *testing.T) {
	d := &ServiceDiscovery{
		cancel: func() {},
		done:   make(chan struct{}),
	}
	returned := make(chan struct{})
	go func() {
		d.Close()
		close(returned)
	}()

	select {
	case <-returned:
		t.Fatal("Close returned before done was closed")
	case <-time.After(50 * time.Millisecond):
	}

	close(d.done)
	select {
	case <-returned:
	case <-time.After(time.Second):
		t.Fatal("Close did not return after done was closed")
	}
}

func TestEtcdDiscoveryIntegration(t *testing.T) {
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

	service := "etcd-discovery-test"
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

	if got := discovery.Instances()["instance-1"]; got != "127.0.0.1:9001" {
		t.Fatalf("initial instance = %q", got)
	}

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-2"), "127.0.0.1:9002"); err != nil {
		t.Fatal(err)
	}
	waitForDiscovery(t, discovery, func(instances map[string]string) bool {
		return instances["instance-2"] == "127.0.0.1:9002"
	})

	if _, err := cli.Put(ctx, ServiceKey(service, "instance-1"), "127.0.0.1:9191"); err != nil {
		t.Fatal(err)
	}
	waitForDiscovery(t, discovery, func(instances map[string]string) bool {
		return instances["instance-1"] == "127.0.0.1:9191"
	})

	deleteResp, err := cli.Delete(ctx, ServiceKey(service, "instance-2"))
	if err != nil {
		t.Fatal(err)
	}
	waitForDiscovery(t, discovery, func(instances map[string]string) bool {
		_, exists := instances["instance-2"]
		return !exists
	})
	if got := discovery.Revision(); got != deleteResp.Header.Revision {
		t.Fatalf("Revision() = %d, want %d", got, deleteResp.Header.Revision)
	}

	instances := discovery.Instances()
	instances["instance-1"] = "mutated"
	if got := discovery.Instances()["instance-1"]; got != "127.0.0.1:9191" {
		t.Fatalf("Instances() exposed internal map, got %q", got)
	}
}

func waitForDiscovery(t *testing.T, d *ServiceDiscovery, condition func(map[string]string) bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if condition(d.Instances()) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("discovery state did not converge: %v", d.Instances())
}
