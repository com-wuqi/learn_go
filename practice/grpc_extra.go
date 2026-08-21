package practice

// ============================================================
// 第三阶段 3.2-S：gRPC 补充学习（status / 重试）
// ============================================================

// calcCheckServerInterceptor 服务端拦截器：当 AddRequest.A+B > 100 时返回
// status.Error(codes.InvalidArgument, "sum too large")，否则交给 handler。
// 进阶：用 status.New + WithDetails 附加一个 calc.AddReply 作为错误详情，再 st.Err()。
func calcCheckServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	return handler(ctx, req)
}

// RunStatusPractice 起带 calcCheckServerInterceptor 的 server，客户端分别用
// A+B=3（成功）和 A+B=101（失败）调用 Add；失败时用 status.FromError 打印 code 和 message。
// 完整流程参考 review/grpc_status.go 的 RunStatusDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure、
// google.golang.org/grpc/codes、google.golang.org/grpc/status。
func RunStatusPractice() {
	// TODO
}

// flakyCalcServerInterceptor 前 2 次 Add 调用返回 Unavailable，之后放行。
// 提示：用 struct 字段 + sync.Mutex 计数，注意只对 "/calc.Calculator/Add" 生效；
// 别用包级裸变量，并发调用会 data race。
type flakyCalcServerInterceptor struct{}

func (f *flakyCalcServerInterceptor) unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	return handler(ctx, req)
}

// RunRetryPractice 起带 flakyCalcServerInterceptor 的 server：
// 1) 不带重试配置调 Add，观察 Unavailable 失败；
// 2) 带 retryPolicy（maxAttempts=4、retryableStatusCodes=["UNAVAILABLE"]）再调，观察自动重试成功。
// 完整流程参考 review/grpc_retry.go 的 RunRetryDemo。
//
// 实现本函数时需要补充 import：context、fmt、net、time、sync、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure、
// google.golang.org/grpc/codes、google.golang.org/grpc/status。
func RunRetryPractice() {
	// TODO
}
