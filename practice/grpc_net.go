package practice

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	return nil
}

// RunTlsPractice 生成 CA/证书并跑四种场景（完整流程参考 review/grpc_tls.go 的 RunTlsDemo）：
// 1) TLS 信任 CA 成功；2) 客户端不信任 CA 失败；
// 3) mTLS 带客户端证书成功；4) mTLS 无客户端证书失败。
// 实现本函数时需要补充 import：fmt、net、time、crypto/tls、calc、insecure（如需）。
func RunTlsPractice() {
	// TODO
}

// watchCalcConnState 打印连接状态变化，直到 ctx 结束。
// 提示：state := conn.GetState() 打印，再 conn.WaitForStateChange(ctx, state) 等下一次变化。
func watchCalcConnState(conn *grpc.ClientConn, ctx context.Context, label string) {
	// TODO
}

// RunConnStatePractice 起 server、建 client（不 Connect），观察状态机：
// 初始 IDLE → Connect 后 CONNECTING/READY → 停 server 变 TRANSIENT_FAILURE → 同地址重启恢复 READY。
// 完整流程参考 review/grpc_connstate.go 的 RunConnStateDemo。
// 实现本函数时需要补充 import：fmt、net、time、calc、insecure。
func RunConnStatePractice() {
	// TODO
}
