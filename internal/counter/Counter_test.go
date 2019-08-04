package counter

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkGetCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCounter(TOKEN_INSERT)
	}
}
func TestStandardCounter_Count(t *testing.T) {
	counter := GetCounter(TOKEN_INSERT)
	waitGroup := new(sync.WaitGroup)
	var loopNum int64
	loopNum = 1000000
	waitGroup.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 10000; i++ {
				counter.Inc(1)
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()

	if counter.Count() == loopNum {
		t.Log("count right", loopNum)
	} else {
		t.Error("count error", loopNum, counter.Count())
	}
}

func BenchmarkStandardCounter_Inc(b *testing.B) {
	counter := GetCounter(TOKEN_INSERT)
	for i := 0; i < b.N; i++ {
		counter.Inc(1)
	}
	fmt.Println(counter.Count())
}

func BenchmarkStandardCounter_Dec(b *testing.B) {
	counter := GetCounter(TOKEN_INSERT)
	for i := 0; i < b.N; i++ {
		counter.Dec(1)
	}
	fmt.Println(counter.Count())
}

func BenchmarkStandardCounter_GetAndInc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCounter(TOKEN_INSERT).Inc(1)
	}
}
