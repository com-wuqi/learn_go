package practice

import (
	"testing"
)

// ============================================================
// 练习 34：表驱动测试 — 为 KVStore 写测试
// 运行: go test -v -run TestKVStore ./practice/
// ============================================================

// [TODO] 用表驱动测试覆盖 MemStore 的 Get/Set/Delete
// tests 表包含多个 test case（name, key, value, want, wantErr）
// 使用 t.Run() 执行子测试
func TestKVStore(t *testing.T) {
	// [TODO] 构造 tests 表，至少覆盖：
	//   1. Set 后 Get 能取到（正常流程）
	//   2. Get 不存在的 key 返回 error
	//   3. Delete 后 Get 返回 error
	//
	// 参考模板:
	// tests := []struct {
	//     name    string
	//     key     string
	//     want    string
	//     wantErr error
	// }{ ... }
	// for _, tt := range tests {
	//     t.Run(tt.name, func(t *testing.T) {
	//         store := NewMemStore()
	//         ...
	//         if got != tt.want { t.Errorf(...) }
	//     })
	// }
}

// ============================================================
// 练习 35：基准测试 — RingBuffer Push/Pop
// 运行: go test -bench=BenchmarkRingBuffer ./practice/
// ============================================================

// [TODO] 对 RingBuffer 的 Push 操作做基准测试
// b.N 是框架自动调整的循环次数
func BenchmarkRingBufferPush(b *testing.B) {
	// [TODO] 用 b.N 做循环次数
}

// [TODO] 对 Push+Pop 组合做基准测试
func BenchmarkRingBufferPushPop(b *testing.B) {
	// [TODO]
}
