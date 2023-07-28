# Go-Day7

## 通道缓冲

默认情况下，通道是无缓冲的，这意味着只有对应的接收（<-chan)，通道准备好接收时，才允许进行发送（chan<-）。有缓冲通道允许在没有对应接收者的情况下，缓冲一定数量的值。

```Go
package main

import "fmt"

func main() {

    messages := make(chan string, 2)

    messages <- "buffered"
    messages <- "channel"

    fmt.Println(<-messages)
    fmt.Println(<-messages)
}
```

在代码中，make了一个字符串通道，最多允许缓存两个值。

由于此通道是有缓冲的，因此我们可以将这些值发送到通道中，而无需并发的接收，然后即可正常接收这两个值。

```shell
$ go run channel-buffering.go 
buffered
channel
```

## 通道同步

我们可以使用通道来同步协程之间的执行状态。这儿有一个例子，使用阻塞接收的方式，实现了等待另一个协程完成。如果需要等待多个协程，WaitGroup是一个更好的选择。

```Go
package main

import (
    "fmt"
    "time"
)

func worker(done chan bool) {
    fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")

    done <- true
}

func main() {

    done := make(chan bool, 1)
    go worker(done)

    <-done
}
```

在代码中，将要在协程中运行这个函数。done通道将被用于通知其他协程这个函数已经完成工作。发送一个值来通知已经完工啦。

运行一个worker协程，并给予用于通知的通道，程序将一直阻塞，直至收到worker使用通道发送的通知。

如果把 <-done这句代码从中删除，程序甚至可能在worker开始运行前就结束了。

```shell
$ go run channel-synchronization.go
working...done
```

## 通道方向

当使用通道作为函数的参数时，可以指定这个通道是否为只读或只写。该特性可以提升程序的类型安全。

代码中，ping函数定义了一个只能发送数据的（只写）通道。尝试从这个通道接收数据会是一个编译时错误。

pong函数接收两个通道，pings仅用于接收数据（只读），pongs仅用于发送数据（只写）。

```Go
package main

import "fmt"

func ping(pings chan<- string, msg string) {
    pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
    msg := <-pings
    pongs <- msg
}

func main() {
    pings := make(chan string, 1)
    pongs := make(chan string, 1)
    ping(pings, "passed message")
    pong(pings, pongs)
    fmt.Println(<-pongs)
}
```

```shell
$ go run channel-directions.go
passed message
```

## 通道选择器

Go的选择器（select），让你可以同时等待多个通道操作，将协程、通道和选择器结合，是Go的一个强大特性。

在代码中，我们将从两个通道中选择

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }
}
```

各个通道将在一定时间后接收一个值，通过这种方式来模拟并行的协程执行（例如，RPC操作）时造成的阻塞（耗时）

使用select关键字来同时等待这两个值，并打印各自接收到的值。跟预期的一样，首先接收到的是值”one“，然后是”two“。

注意，程序总共仅运行了两秒左右。因为 1 秒 和 2 秒的 `Sleeps` 是并发执行的

```Shell
$ time go run select.go 
received one
received two

real    0m2.245s
```

