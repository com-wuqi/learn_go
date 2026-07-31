package practice

import (
	"math/rand"
	"sync"
	"time"
)

// ============================================================
// goroutine + channel 练习（学完 review/goroutine.go 后做）
// ============================================================

// 练习 45：生产者-消费者
// [TODO] 启动 produceCount 个生产者（每个发 1 个数），1 个消费者接收并求和
// produceCount 个 goroutine 各自向 ch 发一个随机数，全部发完后 close(ch)
// 消费者用 range 接收并累加，返回总和
// 提示: sync.WaitGroup 等所有生产者完成后再 close(ch)
func ProducerConsumerSum(produceCount int) int {
	ch1 := make(chan int, produceCount+1)
	wg := &sync.WaitGroup{}
	for range produceCount {
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			ch1 <- rand.Int()
		}(wg)
	}
	wg.Wait()
	close(ch1)
	var sum int
	for v := range ch1 {
		sum += v
	}
	return sum
}

// 练习 14（激活）：交替打印 1-10
// 见 practice/concurrency.go: AlternatePrint
// 提示: 两个 goroutine + 两个 channel 互相通知

// 练习 46：用 select 实现超时后的默认值
// [TODO] SelectWithDefault 从 ch 读取，如果 100ms 内没收到数据返回 defaultVal
// 提示: select + time.After + default
func SelectWithDefault(ch chan int, defaultVal int) int {
	var result int
	select {
	case d, ok := <-ch:
		if ok {
			result = d
		} else {
			result = defaultVal
		}
	case <-time.After(time.Millisecond * 100):
		result = defaultVal
	}
	return result
}
