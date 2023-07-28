package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var ops uint64
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1) // 为每个协程增加一个计数器

		go func() {
			for c := 0; c < 1000; c++ {
				atomic.AddUint64(&ops, 1)
			}
			wg.Done() // 每个协程完成后减少一个计数器
		}()
	}
	wg.Wait() // 等待所有协程完成

	fmt.Println("ops:", ops)
}
