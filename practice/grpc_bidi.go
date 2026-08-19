package practice

import (
	"LearnGo/api/calc"
)

// ============================================================
// 第三阶段 3.2：gRPC Bidirectional-streaming RPC
// ============================================================

// EchoSum 实现双向流：循环 Recv，把每个 AddRequest 的 a+b 逐条回发。
// 提示：
//   - stream 类型是 calc.Calculator_EchoSumServer
//   - 服务端 Recv 到 io.EOF（客户端 CloseSend）时 return nil
//   - 每条消息用 stream.Send(&calc.AddReply{Result: req.A + req.B}) 回发
func (s *calculatorServer) EchoSum(stream calc.Calculator_EchoSumServer) error {
	// TODO
	return nil
}

// RunCalculatorBidiDemo 起 server、建 client、边发边收，最后 CloseSend。
// 完整流程参考 review/grpc_bidi.go 的 RunBidiDemo。
//
// 实现本函数时需要补充 import：context、fmt、io、net、time、
// google.golang.org/grpc、google.golang.org/grpc/credentials/insecure。
func RunCalculatorBidiDemo() {
	// TODO
}
