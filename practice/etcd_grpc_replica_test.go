package practice

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"LearnGo/api/demo"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestEtcdRegisteredGreeterLifecycle(t *testing.T) {
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

	service := "etcd-registered-greeter-test"
	instanceID := "instance-1"
	key := ServiceKey(service, instanceID)
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cleanupCancel()
		_, _ = cli.Delete(cleanupCtx, ServicePrefix(service), clientv3.WithPrefix())
	}()

	replica, err := StartRegisteredGreeter(ctx, cli, service, instanceID, "replica-A", 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := replica.Stop(); err != nil {
			t.Errorf("Stop() = %v", err)
		}
	})

	if replica.Addr == "" {
		t.Fatal("replica Address is empty")
	}
	getResp, err := cli.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	if len(getResp.Kvs) != 1 || string(getResp.Kvs[0].Value) != replica.Addr {
		t.Fatalf("registered value = %#v, want %q", getResp.Kvs, replica.Addr)
	}

	conn, err := grpc.NewClient(replica.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	callCtx, callCancel := context.WithTimeout(ctx, 2*time.Second)
	defer callCancel()
	reply, err := demo.NewGreeterClient(conn).SayHello(callCtx, &demo.HelloRequest{Name: "Codex"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reply.GetMessage(), "replica-A") {
		t.Fatalf("response = %q, want replica label", reply.GetMessage())
	}

	if err := replica.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := replica.Stop(); err != nil {
		t.Fatalf("second Stop() = %v", err)
	}

	getResp, err = cli.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	if len(getResp.Kvs) != 0 {
		t.Fatalf("registration remains after Stop: %#v", getResp.Kvs)
	}
}
