package practice

// ============================================================
// 第三阶段 3.2：gRPC metadata / 拦截器（可选进阶）
// ============================================================

// loggingUnaryServerInterceptor 服务端 Unary 拦截器：
// 打印方法名和耗时；若 incoming metadata 带 x-user-id，也打印出来。
// 签名：func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
// 提示：metadata.FromIncomingContext(ctx) 取 metadata；handler(ctx, req) 调用链上的下一个处理器。
func loggingUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	return handler(ctx, req)
}

// loggingUnaryClientInterceptor 客户端 Unary 拦截器：打印方法名和耗时。
// 签名：func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
// 提示：invoker(ctx, method, req, reply, cc, opts...) 发起真实调用。
func loggingUnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// TODO
	return invoker(ctx, method, req, reply, cc, opts...)
}

// RunCalculatorAdvancedDemo 起 server（带服务端拦截器）、建 client（带客户端拦截器），
// 用 metadata.AppendToOutgoingContext 附加 x-user-id 后调用 Add，观察两侧日志。
// 完整流程参考 review/grpc_advanced.go 的 RunMetadataDemo / RunInterceptorDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、
// google.golang.org/grpc/metadata、google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorAdvancedDemo() {
	// TODO
}
