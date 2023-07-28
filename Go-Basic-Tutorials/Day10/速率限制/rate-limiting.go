package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5) // 通道缓冲区大小为5
	for i := 1; i <= 5; i++ {
		requests <- i
	}

	close(requests) // 关闭通道

	// 限制速率
	limiter := time.Tick(time.Millisecond * 200)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	burstyLimiter := make(chan time.Time, 3) // 通道缓冲区大小为3
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Millisecond * 200) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}

	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}
