# Go-Day9

## 定时器-Timer

我么经常需要在未来的某个时间点运行Go代码，或者每隔一段时间重复运行代码，Go内置的定时器和打点器特性让这些操作变得十分简单。

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    timer1 := time.NewTimer(2 * time.Second)

    <-timer1.C
    fmt.Println("Timer 1 fired")

    timer2 := time.NewTimer(time.Second)
    go func() {
        <-timer2.C
        fmt.Println("Timer 2 fired")
    }()
    stop2 := timer2.Stop()
    if stop2 {
        fmt.Println("Timer 2 stopped")
    }

    time.Sleep(2 * time.Second)
}
```

定时器表示在未来某一时刻的独立事件，你告诉定时器需要等待的时间，然后它将提供一个用于通知的通道，代码中的定时器会等待2秒

<-timer1.c会一直阻塞，直到定时器的通道c明确的发送了定时器失效的值

如果需要的仅仅是单纯的等待，使用time.Sleep就够了，使用定时器的原因之一就是，可以在定时器触发之前就将其取消掉。

给timer2足够的时间来触发它，以证明它实际上已经停止了。

第一个定时器将在程序开始后大约2s触发，但是第二个定时器还未触发就停止了。

```shell
$ go run timers.go
Timer 1 fired
Timer 2 stopped
```

## 打点器-Ticker

定时器是当你想要在未来某一刻执行一次时使用的，打点器则是为你想要以固定的时间间隔重复执行而准备的。代码展示了一个打点器的例子，它将定时的执行，直到我们将它停止。

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan bool)

    go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                fmt.Println("Tick at", t)
            }
        }
    }()

    time.Sleep(1600 * time.Millisecond)
    ticker.Stop()
    done <- true
    fmt.Println("Ticker stopped")
}
```

打点器和定时器的机制有点相似：使用一个通道来发送数据，这里使用通道内建的select，等待每500ms到达一次的值。

打点器可以和定时器一样被停止，打点器一旦停下，将不能再从它的通道中接收到值。我们将在运行1600ms后停止这个打点器。

当运行程序时，打点器会在我们停止它前打点3次。

```shell
$ go run tickers.go
Tick at 2012-09-23 11:29:56.487625 -0700 PDT
Tick at 2012-09-23 11:29:56.988063 -0700 PDT
Tick at 2012-09-23 11:29:57.488076 -0700 PDT
Ticker stopped
```

## 工作池

在代码中，展示了如何使用协程与通道实现一个工作池。

```Go
package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("worker", id, "started  job", j)
        time.Sleep(time.Second)
        fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}

func main() {

    const numJobs = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= numJobs; a++ {
        <-results
    }
}
```

worker程序，会并发的运行多个worker，worker将在jobs频道上接收工作，并在results上发送相应的结果，每个worker我们都会sleep一秒钟，以模拟一项昂贵的（耗时一秒钟的）任务。

为了使用worker工作池并且收集其的结果，我们需要2个通道。

这里启动了3个worker，初始是阻塞的，因为还没有传递任务。

这里发送了5个jobs，然后close这些通道，表示这些就是所有的任务了。

最后，收集所有这些任务的返回值，也确保了所有的worker的协程都已完成，另一个等待多个协程的方法是使用 `WaitGroup`

运行程序，显示5个任务被多个worker执行，尽管所有的工作总共要花费5秒钟，但该程序只花了2秒钟，因为3个worker是并行的。

```shell
$ time go run worker-pools.go 
worker 1 started  job 1
worker 2 started  job 2
worker 3 started  job 3
worker 1 finished job 1
worker 1 started  job 4
worker 2 finished job 2
worker 2 started  job 5
worker 3 finished job 3
worker 1 finished job 4
worker 2 finished job 5
real    0m2.358s
```

## WaitGroup

想要等待多个协程完成，我们可以使用WaitGroup

每个协程都会运行worker函数，睡眠一秒钟以模拟耗时的任务

```Go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int) {
    fmt.Printf("Worker %d starting\n", id)

    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {

    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)

        i := i

        go func() {
            defer wg.Done()
            worker(i)
        }()
    }

    wg.Wait()

}
```

WaitGroup用于等待这里启动的所有协程完成，注意：如果WaitGroup显式传递到函数中，则应使用指针。

启动几个协程，并为其递增WaitGroup的计数器。

避免在每个协程闭包中重复利用相同的i值可以参考[the FAQ](https://golang.org/doc/faq#closures_and_goroutines)

将worker调用包装在一个闭包中，可以确保通知WaitGroup此工作线程已完成。这样，worker线程本身就不必直到其执行中涉及的并发原语。

阻塞，直到WaitGroup计数器恢复为0，即所有协程的工作都已经完成。

请注意，WaitGroup的使用方式并没有直观的办法传递来自worker的错误，更高级的示例，参考[errgroup package](https://pkg.go.dev/golang.org/x/sync/errgroup)

```Go
$ go run waitgroups.go
Worker 5 starting
Worker 3 starting
Worker 4 starting
Worker 1 starting
Worker 2 starting
Worker 4 done
Worker 1 done
Worker 2 done
Worker 5 done
Worker 3 done
```

每次运行，各个协程开启和完成的时间可能是不同的。