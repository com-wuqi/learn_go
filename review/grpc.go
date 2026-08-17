package review

import (
	"fmt"

	"LearnGo/api/demo"

	"google.golang.org/protobuf/proto"
)

// RunProtoDemo 演示 protoc 生成的消息代码能做什么：
// 构造消息 -> 序列化成二进制 -> 反序列化回来。
// 这是 3.2-A 的演示，真正的 gRPC server/client 在 3.2-B 再写。
func RunProtoDemo() {
	fmt.Println("\n=== proto 生成代码演示：消息序列化 ===")

	req := &demo.HelloRequest{
		Name:    "Codex",
		Age:     18,
		IsVip:   true,
		Tags:    []string{"go", "grpc"},
		Level:   demo.Level_LEVEL_PREMIUM,
		Address: &demo.Address{City: "Shanghai", Street: "Zhangyang"},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		fmt.Println("  marshal 失败:", err)
		return
	}
	fmt.Printf("  序列化后 %d 字节：%x\n", len(data), data)

	var out demo.HelloRequest
	if err := proto.Unmarshal(data, &out); err != nil {
		fmt.Println("  unmarshal 失败:", err)
		return
	}
	fmt.Printf("  反序列化结果：%+v\n", &out)
}
