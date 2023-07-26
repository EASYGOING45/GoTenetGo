# Go-Day10

## 速率限制

速率限制是控制服务资源利用和质量的重要机制。基于协程、通道和打点器，Go优雅的支持速率限制。

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    requests := make(chan int, 5)
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    limiter := time.Tick(200 * time.Millisecond)

    for req := range requests {
        <-limiter
        fmt.Println("request", req, time.Now())
    }

    burstyLimiter := make(chan time.Time, 3)

    for i := 0; i < 3; i++ {
        burstyLimiter <- time.Now()
    }

    go func() {
        for t := range time.Tick(200 * time.Millisecond) {
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
```

首先，看到一个最基本的速率限制，假设我们想限制对收到请求的处理，可以通过一个渠道处理这些请求。

`limiter`通道每200ms接收一个值，这是我们任务速率限制的调度器。

通过在每次请求前阻塞limiter通道的一个接收，可以将频率限制为，每200ms执行一次请求。

有时候我们可能希望在速率限制方案中允许短暂的并发请求，并同事保留总体速率限制，可以通过缓冲通道来完成此任务，burstyLimiter通道允许最多3个爆发（bursts）事件。

填充通道，表示允许的爆发（bursts）。

每200ms我们将尝试添加一个新的值到burstyLimiter中，直到达到3个的限制。

随后，模拟另外5个传入请求。受益于BurstyLimiter的爆发（bursts）能力，前3个请求可以快速完成。

运行程序后，我们可以看到第一批请求意料之中的大约每200ms处理一次。

第二批请求，由于爆发（burstable）速率控制，我们直接连续处理了3个请求，然后以大约每200ms一次的速度，处理了剩余的2个请求。

```shell
$ go run rate-limiting.go
request 1 2012-10-19 00:38:18.687438 +0000 UTC
request 2 2012-10-19 00:38:18.887471 +0000 UTC
request 3 2012-10-19 00:38:19.087238 +0000 UTC
request 4 2012-10-19 00:38:19.287338 +0000 UTC
request 5 2012-10-19 00:38:19.487331 +0000 UTC

request 1 2012-10-19 00:38:20.487578 +0000 UTC
request 2 2012-10-19 00:38:20.487645 +0000 UTC
request 3 2012-10-19 00:38:20.487676 +0000 UTC
request 4 2012-10-19 00:38:20.687483 +0000 UTC
request 5 2012-10-19 00:38:20.887542 +0000 UTC
```

## 原子计数器

Go中最主要的状态管理机制是依靠通道间的通信来完成的。之前在工作池的例子中看到过，并且，还有一些其他的方法来管理状态，在代码中，展示了如何使用 `sync/atomic`包在多个协程中进行 __原子计数__。

```Go
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
        wg.Add(1)

        go func() {
            for c := 0; c < 1000; c++ {

                atomic.AddUint64(&ops, 1)
            }
            wg.Done()
        }()
    }

    wg.Wait()

    fmt.Println("ops:", ops)
}
```

使用一个无符号整型（永远是正整数）变量来表示这个计数器

WaitGroup帮助我们等待所有协程完成它们的工作

启动50个协程，并且每个协程会将计数器递增1000次。

使用AddUint64来让计数器自动增加，使用&语法给定ops的内存地址。

```Go
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
        wg.Add(1)

        go func() {
            for c := 0; c < 1000; c++ {

                atomic.AddUint64(&ops, 1)
            }
            wg.Done()
        }()
    }

    wg.Wait()

    fmt.Println("ops:", ops)
}
```

随后，等待，直到所有协程完成

现在可以安全的访问ops，因为我们知道，此时没有协程写入ops，此外，还可以使用 atomic.LoadUint64之类的函数，在原子更新的同时安全地读取它们。

预计会进行50000次操作，如果使用非原子的ops++来增加计数器，由于多个协程会互相干扰，运行时值会改变，可能会导致我们得到一个不同的数字，此外，运行程序时带上 -race标志，我们可以获取数据竞争失败的详情。

```shell
$ go run atomic-counters.go
ops: 50000
```

## 互斥锁

之前讲了如何使用原子操作（atomic-counters）来管理简单的计数器，对于更复杂的情况，可以使用一个互斥量来在Go协程间安全的访问数据。

```Go
package main

import (
    "fmt"
    "sync"
)

type Container struct {
    mu       sync.Mutex
    counters map[string]int
}

func (c *Container) inc(name string) {

    c.mu.Lock()
    defer c.mu.Unlock()
    c.counters[name]++
}

func main() {
    c := Container{

        counters: map[string]int{"a": 0, "b": 0},
    }

    var wg sync.WaitGroup

    doIncrement := func(name string, n int) {
        for i := 0; i < n; i++ {
            c.inc(name)
        }
        wg.Done()
    }

    wg.Add(3)
    go doIncrement("a", 10000)
    go doIncrement("a", 10000)
    go doIncrement("b", 10000)

    wg.Wait()
    fmt.Println(c.counters)
}
```

Container中定义了counters的map，由于我们希望从多个goroutine同时更新它，因此我们添加了一个互斥锁Mutex来同步访问，请注意，不能复制互斥锁，如果需要传递这个struct，应该使用指针完成。

在访问counters之前锁定互斥锁，使用`[defer](defer)`在函数结束时解锁。

注意，互斥量的零值是可用的，因此这里不需要初始化。

doIncrement函数在循环中递增对name的计数。

同时运行多个goroutines，注意，它们都访问相同的Container，其中两个访问相同的计数器。

```shell
$ go run mutexes.go
map[a:20000 b:10000]
```

## 状态协程

在之前，讲了如何用互斥锁进行明确的锁定，来让共享的state跨多个Go协程同步访问，另一个选择是，使用内建协程和通道的同步特性来达到同样的效果，Go共享内存的思想史：通过通信使每个数据仅被单个协程所拥有，即通过通信实现共享内存，就通道的方法与该思想完全一致！

```Go
package main

import (
    "fmt"
    "math/rand"
    "sync/atomic"
    "time"
)

type readOp struct {
    key  int
    resp chan int
}
type writeOp struct {
    key  int
    val  int
    resp chan bool
}

func main() {

    var readOps uint64
    var writeOps uint64

    reads := make(chan readOp)
    writes := make(chan writeOp)

    go func() {
        var state = make(map[int]int)
        for {
            select {
            case read := <-reads:
                read.resp <- state[read.key]
            case write := <-writes:
                state[write.key] = write.val
                write.resp <- true
            }
        }
    }()

    for r := 0; r < 100; r++ {
        go func() {
            for {
                read := readOp{
                    key:  rand.Intn(5),
                    resp: make(chan int)}
                reads <- read
                <-read.resp
                atomic.AddUint64(&readOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    for w := 0; w < 10; w++ {
        go func() {
            for {
                write := writeOp{
                    key:  rand.Intn(5),
                    val:  rand.Intn(100),
                    resp: make(chan bool)}
                writes <- write
                <-write.resp
                atomic.AddUint64(&writeOps, 1)
                time.Sleep(time.Millisecond)
            }
        }()
    }

    time.Sleep(time.Second)

    readOpsFinal := atomic.LoadUint64(&readOps)
    fmt.Println("readOps:", readOpsFinal)
    writeOpsFinal := atomic.LoadUint64(&writeOps)
    fmt.Println("writeOps:", writeOpsFinal)
}
```

在代码中，state将被一个单独的协程拥有，这能保证数据在并行读取时不会混乱，为了对state进行读取或者写入，其他的协程将发送一条数据到目前拥有数据的协程中，然后等待接收对应的恢复。结构体readOp和writeOp封装了这些请求，并提供了相应协程的方法。

和前面一样，我们会计算操作执行的次数。

其他协程将通过reads和writes通道来发布 读 和 写请求。

这就是拥有state的那个协程，和前面例子中的map一样，不过这里的state是被这个状态协程私有的，这个协程不断地在reads和writes通道上进行选择，并在请求到达时做出相应。首先，执行请求的操作；然后，执行响应，在响应通道resp上发送一个值，表明请求成功（reads的值则为state对应的值）

启动100个协程通过reads通道向拥有state的协程发起读取请求，每个读取请求需要构造一个readOp，发送它到reads通道中，并通过给定的resp通道接收结果。

用相同的方法启动10个写操作。

让协程们跑1s，最后，获取并报告ops值。

运行这个程序后显示这个基于协程的状态管理的例子，达到了每秒大约80000次操作。

通过这个例子我们可以看到，基于协程的方法比基于互斥锁的方法要复杂得多，但是，在某些情况下它可能很有用，例如，当你涉及其他通道，或者管理多个同类互斥锁时，会很容易出错。应该使用最自然的方法，尤其是在理解程序正确性方面

```shell
$ go run stateful-goroutines.go
readOps: 71708
writeOps: 7177
```

