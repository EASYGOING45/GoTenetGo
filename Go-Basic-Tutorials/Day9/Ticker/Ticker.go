package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(500 * time.Millisecond) // 创建一个定时器，每500毫秒触发一次
	done := make(chan bool)                          // 创建一个通道，用于接收定时器的信号

	go func() {
		for {
			select {
			case <-done: // 通道接收到信号
				return
			case t := <-ticker.C: // 定时器触发
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond) // 等待1600毫秒
	ticker.Stop()                       // 停止定时器
	done <- true                        // 发送信号

	fmt.Println("Ticker stopped")
}
