package practice

import (
	"LearnGo/api/calc"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ============================================================
// 第三阶段 3.2-S 批2：TLS / 连接状态机
// ============================================================

// ---- 证书工具（样板代码，理解即可，不用改）----

// newSelfSignedCA 生成自签名 CA（证书 + 私钥）。
func newSelfSignedCA() (*x509.Certificate, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "LearnGo Practice CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		panic(err)
	}
	return cert, key
}

// issueCert 用 CA 签发 server/client 证书（含 localhost/127.0.0.1 SAN）。
func issueCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, cn string, isServer bool) tls.Certificate {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}
	if isServer {
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
	if err != nil {
		panic(err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// newCAPool 把 CA 证书放进信任池。
func newCAPool(caCert *x509.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AddCert(caCert)
	return pool
}

// ---- 练习 TODO ----

// tlsClientCreds 返回客户端 TLS credentials：
// caPool 用于校验服务端证书；useMTLS 为 true 时带上 clientCert；
// serverName 是证书校验用的主机名。
// 提示：credentials.NewTLS(&tls.Config{RootCAs: caPool, Certificates: ..., ServerName: serverName})
func tlsClientCreds(caPool *x509.CertPool, clientCert tls.Certificate, serverName string, useMTLS bool) credentials.TransportCredentials {
	// TODO
	baseConfig := &tls.Config{RootCAs: caPool, ServerName: serverName}
	if useMTLS {
		baseConfig.Certificates = []tls.Certificate{clientCert}
	}
	return credentials.NewTLS(baseConfig)
}

// tlsServerCreds 返回服务端 TLS credentials：
// serverCert 是服务端证书；requireClientCert 为 true 时启用 mTLS
// （ClientAuth: tls.RequireAndVerifyClientCert + ClientCAs: caPool）。
func tlsServerCreds(serverCert tls.Certificate, caPool *x509.CertPool, requireClientCert bool) credentials.TransportCredentials {
	// TODO
	baseConfig := &tls.Config{Certificates: []tls.Certificate{serverCert}}
	if requireClientCert {
		baseConfig.ClientAuth = tls.RequireAndVerifyClientCert
		baseConfig.ClientCAs = caPool
	}
	return credentials.NewTLS(baseConfig)
}

// RunTlsPractice 生成 CA/证书并跑四种场景（完整流程参考 review/grpc_tls.go 的 RunTlsDemo）：
// 1) TLS 信任 CA 成功；2) 客户端不信任 CA 失败；
// 3) mTLS 带客户端证书成功；4) mTLS 无客户端证书失败。
// 实现本函数时需要补充 import：fmt、net、time、crypto/tls、calc、insecure（如需）。
func RunTlsPractice() {
	// TODO
	cert, key := newSelfSignedCA()
	serverCert := issueCert(cert, key, "server", true)
	clientCert := issueCert(cert, key, "client", false)
	caPool := newCAPool(cert)

	ln1, err := net.Listen("tcp", "127.0.0.1:0")
	// TLS
	if err != nil {
		fmt.Println(err)
		return
	}

	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	// mTLS
	if err != nil {
		fmt.Println(err)
		return
	}

	server1 := grpc.NewServer(grpc.Creds(tlsServerCreds(serverCert, caPool, false)))
	// TLS
	calc.RegisterCalculatorServer(server1, &calculatorServer{})
	go func() {
		err := server1.Serve(ln1)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer server1.GracefulStop()

	server2 := grpc.NewServer(grpc.Creds(tlsServerCreds(serverCert, caPool, true)))
	// mTLS
	calc.RegisterCalculatorServer(server2, &calculatorServer{})
	go func() {
		err := server2.Serve(ln2)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer server2.GracefulStop()

	conn1, err := grpc.NewClient(
		ln1.Addr().String(),
		grpc.WithTransportCredentials(tlsClientCreds(caPool, tls.Certificate{}, "127.0.0.1", false)),
		// TLS-success
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn1 *grpc.ClientConn) {
		err := conn1.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn1)
	client1 := calc.NewCalculatorClient(conn1)
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	resp1, err := client1.Add(ctx1, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp1.GetResult())
	}

	conn2, err := grpc.NewClient(
		ln1.Addr().String(),
		grpc.WithTransportCredentials(tlsClientCreds(x509.NewCertPool(), tls.Certificate{}, "127.0.0.1", false)),
		// TLS-fail
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn2 *grpc.ClientConn) {
		err := conn2.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn2)
	client2 := calc.NewCalculatorClient(conn2)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	resp2, err := client2.Add(ctx2, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.GetResult())
	}

	conn3, err := grpc.NewClient(
		ln2.Addr().String(),
		grpc.WithTransportCredentials(tlsClientCreds(caPool, clientCert, "127.0.0.1", true)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn3 *grpc.ClientConn) {
		err := conn3.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn3)
	client3 := calc.NewCalculatorClient(conn3)
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()
	resp3, err := client3.Add(ctx3, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp3.GetResult())
	}

	conn4, err := grpc.NewClient(
		ln2.Addr().String(),
		grpc.WithTransportCredentials(tlsClientCreds(caPool, tls.Certificate{}, "127.0.0.1", false)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn4 *grpc.ClientConn) {
		err := conn4.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn4)
	client4 := calc.NewCalculatorClient(conn4)
	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel4()
	resp4, err := client4.Add(ctx4, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp4.GetResult())
	}

}

// watchCalcConnState 打印连接状态变化，直到 ctx 结束。
// 提示：state := conn.GetState() 打印，再 conn.WaitForStateChange(ctx, state) 等下一次变化。
func watchCalcConnState(conn *grpc.ClientConn, ctx context.Context, label string) {
	// TODO
	for {
		state := conn.GetState()
		fmt.Printf("[%s] state=%s\n", label, state)
		if !conn.WaitForStateChange(ctx, state) {
			return
		}
	}
}

// RunConnStatePractice 起 server、建 client（不 Connect），观察状态机：
// 初始 IDLE → Connect 后 CONNECTING/READY → 停 server 变 TRANSIENT_FAILURE → 同地址重启恢复 READY。
// 完整流程参考 review/grpc_connstate.go 的 RunConnStateDemo。
// 实现本函数时需要补充 import：fmt、net、time、calc、insecure。
func RunConnStatePractice() {
	// TODO
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	addr := ln.Addr().String()
	srv := grpc.NewServer()
	calc.RegisterCalculatorServer(srv, &calculatorServer{})
	go func() {
		err := srv.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()

	conn, err := grpc.NewClient(ln.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	go watchCalcConnState(conn, ctx, "RunConnStatePractice")

	fmt.Printf("[RunConnStatePractice] state=%s\n", conn.GetState())
	conn.Connect()
	client := calc.NewCalculatorClient(conn)
	firstCallCtx, firstCallCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer firstCallCancel()
	resp, err := client.Add(firstCallCtx, &calc.AddRequest{A: 1, B: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.GetResult())
	}

	srv.GracefulStop()
	failCtx, failCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer failCancel()
	resp2, err := client.Add(failCtx, &calc.AddRequest{A: 1, B: 2}, grpc.WaitForReady(true)) // 等待
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.GetResult())
	}

	ln2, err := net.Listen("tcp", addr)
	for i := 0; err != nil && i < 20; i++ {
		time.Sleep(time.Millisecond * 200)
		ln2, err = net.Listen("tcp", addr)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	srv2 := grpc.NewServer()
	calc.RegisterCalculatorServer(srv2, &calculatorServer{})
	go func() {
		err := srv2.Serve(ln2)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println(err)
		}
	}()
	defer srv2.GracefulStop()
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	resp3, err := client.Add(ctx2, &calc.AddRequest{A: 1, B: 2}, grpc.WaitForReady(true))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp3.GetResult())
	}

}
