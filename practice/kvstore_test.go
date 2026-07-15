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
	tests := []struct {
		name       string
		init       map[string]string
		key        string
		want       string
		needDelete []string
		wantErr    bool
	}{
		{name: "test1：Get存在的key", init: map[string]string{"a": "a"}, key: "a", want: "a", wantErr: false},
		{name: "test2: Get不存在的key", init: map[string]string{"a": "a"}, key: "b", want: "", wantErr: true},
		{name: "test3: Delete存在的key", init: map[string]string{"a": "a"}, key: "a", want: "a", needDelete: []string{"a"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemStore()
			for k, v := range tt.init {
				err := store.Set(k, v)
				if err != nil {
					t.Errorf("Set() error = %v", err)
				}
			}
			if got, err := store.Get(tt.key); (err != nil) != tt.wantErr || got != tt.want {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.needDelete != nil {
				for _, d := range tt.needDelete {
					err := store.Delete(d)
					if err != nil {
						t.Errorf("Delete() error = %v", err)
					}
					_, err = store.Get(d)
					if err == nil {
						t.Errorf("Delete() should have deleted %v", d)
					}
				}
			}

		})
	}
}

// ============================================================
// 练习 35：基准测试 — RingBuffer Push/Pop
// 运行: go test -bench=BenchmarkRingBuffer ./practice/
// ============================================================

// [TODO] 对 RingBuffer 的 Push 操作做基准测试
// b.N 是框架自动调整的循环次数
func BenchmarkRingBufferPush(b *testing.B) {
	// [TODO] 用 b.N 做循环次数
	ring := NewRingBuffer(b.N) // 这里换成一个常数就是 对 Push+Pop 组合做基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ring.IsFull() {
			_, err := ring.Pop()
			if err != nil {
				b.Errorf("Pop() error = %v", err)
			}
		}
		err := ring.Push(i)
		if err != nil {
			b.Errorf("Push() error = %v", err)
		}
	}
	b.StopTimer()
	b.ReportAllocs()

}

//// 对 Push+Pop 组合做基准测试
//func BenchmarkRingBufferPushPop(b *testing.B) {
//	//
//}
