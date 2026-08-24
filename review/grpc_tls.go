package review

import (
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

	"LearnGo/api/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// newSelfSignedCA 生成一个自签名 CA（证书 + 私钥）。
func newSelfSignedCA() (*x509.Certificate, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "LearnGo Demo CA"},
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

// issueCert 用 CA 签发一张 server 或 client 证书（含 localhost/127.0.0.1 SAN）。
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

// caPool 把 CA 证书放入信任池。
func caPool(caCert *x509.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AddCert(caCert)
	return pool
}

// RunTlsDemo 演示 TLS 与 mTLS：四种场景对比。
func RunTlsDemo() {
	fmt.Println("\n=== gRPC TLS / mTLS 演示：证书链与双向认证 ===")

	caCert, caKey := newSelfSignedCA()
	serverCert := issueCert(caCert, caKey, "server", true)
	clientCert := issueCert(caCert, caKey, "client", false)
	pool := caPool(caCert)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 场景 1/2：TLS-only server（只校验客户端是否信任服务端证书）。
	ln1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}
	srv1 := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{serverCert}})))
	demo.RegisterGreeterServer(srv1, &greeterServer{})
	go func() {
		if err := srv1.Serve(ln1); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv1.Stop()

	// 场景 1：客户端信任 CA → 成功。
	connOK, err := grpc.NewClient(ln1.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: pool, ServerName: "127.0.0.1"})))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	resp, err := demo.NewGreeterClient(connOK).SayHello(ctx, &demo.HelloRequest{Name: "Alice"})
	fmt.Println("  1) TLS 信任 CA:", resp.GetMessage(), err)
	connOK.Close()

	// 场景 2：客户端不信任 CA（RootCAs 为空，走系统信任库）→ 失败。
	connBad, err := grpc.NewClient(ln1.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "127.0.0.1"})))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	_, err = demo.NewGreeterClient(connBad).SayHello(ctx, &demo.HelloRequest{Name: "Bob"})
	fmt.Println("  2) TLS 不信任 CA:", err)

	// 场景 3/4：mTLS server（要求并验证客户端证书）。
	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("  监听失败:", err)
		return
	}
	srv2 := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	})))
	demo.RegisterGreeterServer(srv2, &greeterServer{})
	go func() {
		if err := srv2.Serve(ln2); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println("  Serve 异常退出:", err)
		}
	}()
	defer srv2.Stop()

	// 场景 3：客户端带证书 → 成功。
	connMTLS, err := grpc.NewClient(ln2.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs:      pool,
			ServerName:   "127.0.0.1",
			Certificates: []tls.Certificate{clientCert},
		})))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	resp, err = demo.NewGreeterClient(connMTLS).SayHello(ctx, &demo.HelloRequest{Name: "Carol"})
	fmt.Println("  3) mTLS 带客户端证书:", resp.GetMessage(), err)

	// 场景 4：客户端不带证书 → 握手失败。
	connNoCert, err := grpc.NewClient(ln2.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: pool, ServerName: "127.0.0.1"})))
	if err != nil {
		fmt.Println("  建立连接失败:", err)
		return
	}
	_, err = demo.NewGreeterClient(connNoCert).SayHello(ctx, &demo.HelloRequest{Name: "Dave"})
	fmt.Println("  4) mTLS 无客户端证书:", err)
}
