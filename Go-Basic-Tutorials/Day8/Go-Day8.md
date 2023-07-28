# Go-Day8

## 超时处理

超时 对于一个需要连接外部资源，或者有耗时较长的操作的程序而言是很重要的。得益于通道和select，在Go中实现超时操作是简洁而优雅的。

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    c1 := make(chan string, 1)
    go func() {
        time.Sleep(2 * time.Second)
        c1 <- "result 1"
    }()

    select {
    case res := <-c1:
        fmt.Println(res)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout 1")
    }

    c2 := make(chan string, 1)
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "result 2"
    }()
    select {
    case res := <-c2:
        fmt.Println(res)
    case <-time.After(3 * time.Second):
        fmt.Println("timeout 2")
    }
}
```

在代码中，假如我们执行一个外部调用，并在2秒后使用通道c1来返回它的执行结果。

这里是使用select实现一个超时操作。 `res := <- c1`等待结果， `<- time.Ater`等待超时（1秒钟）以后发送的值。由于select默认处理第一个已准备好的接收操作，因此如果操作耗时超过了允许的一秒的话，将会执行超时case。

如果允许一个长一点的超时时间：3秒，就可以成功的从c2接收到值，并且打印出结果。

运行程序后，首先显示运行超时的操作，然后是成功接收的。

```shell
$ go run timeouts.go 
timeout 1
result 2
```

## 非阻塞通道操作

常规的通过通道发送和接收数据是阻塞的。然而，我们可以使用带一个default子句的select来实现非阻塞的发送、接收，甚至是非阻塞的多路select。

```Go
package main

import "fmt"

func main() {
    messages := make(chan string)
    signals := make(chan bool)

    select {
    case msg := <-messages:
        fmt.Println("received message", msg)
    default:
        fmt.Println("no message received")
    }

    msg := "hi"
    select {
    case messages <- msg:
        fmt.Println("sent message", msg)
    default:
        fmt.Println("no message sent")
    }

    select {
    case msg := <-messages:
        fmt.Println("received message", msg)
    case sig := <-signals:
        fmt.Println("received signal", sig)
    default:
        fmt.Println("no activity")
    }
}
```

代码实现了一个非阻塞接收的例子，如果在messages钟存在，然后select将这个值带入 `<-messages case`中。否则，就直接到default分支中。

后续是一个非阻塞发送的例子，代码结构和上面接收的类似。msg不能被发送到message通道，因为这是个无缓冲区通道，并且也没有接收者，因此，执行default。

我们可以在default前使用多个case子句来实现一个多路的非阻塞的选择器，这里我们试图在messages和signals上同时使用非阻塞的接收操作。

```shell
$ go run non-blocking-channel-operations.go 
no message received
no message sent
no activity
```

## 通道的关闭

关闭一个通道意味着不能再向这个通道发送值了。该特性可以向通道的接收方传达工作已经完成的信息。

```Go
package main

import "fmt"

func main() {
    jobs := make(chan int, 5)
    done := make(chan bool)

    go func() {
        for {
            j, more := <-jobs
            if more {
                fmt.Println("received job", j)
            } else {
                fmt.Println("received all jobs")
                done <- true
                return
            }
        }
    }()

    for j := 1; j <= 3; j++ {
        jobs <- j
        fmt.Println("sent job", j)
    }
    close(jobs)
    fmt.Println("sent all jobs")

    <-done
}
```

在例子中，我们将使用一个jobs通道，将工作内容，从main()协程传递到一个工作协程中。当我们没有更多的任务传递给工作协程时，我们将close这个jobs通道。

这是工作协程 go func，使用 `j,more:=<-jobs`循环的从jobs接收数据，根据接收的第二个值，如果jobs已经关闭并且通道中所有的之都已经接收完毕，那么more的值将是false，当我们完成所有的任务时，会使用这个特性通过done通道通知main协程。

使用jobs发送3个任务到工作协程中，然后关闭jobs。

使用前面学习过的通道同步方法等待任务结束。

```shell
$ go run closing-channels.go
sent job 1
received job 1
sent job 2
received job 2
sent job 3
received job 3
sent all jobs
received all jobs
```

根据关闭通道的思想，引出我们的下一个东西，遍历通道。

## 通道遍历

在之前，学习过for和range为基本的数据结构提供了迭代的功能，同样，也可以使用他们来遍历的从通道中取值。

```Go
package main

import "fmt"

func main() {

    queue := make(chan string, 2)
    queue <- "one"
    queue <- "two"
    close(queue)

    for elem := range queue {
        fmt.Println(elem)
    }
}
```

代码便利了queue通道中的两个值

range迭代从queue中得到每个值，因为在前面close了这个通道，所以，这个迭代会在接受完两个值后结束。

代码告诉我们，一个非空的通道也是可以关闭的，并且，通道中剩下的值仍然可以被接收到。

```shell
$ go run range-over-channels.go
one
two
```

