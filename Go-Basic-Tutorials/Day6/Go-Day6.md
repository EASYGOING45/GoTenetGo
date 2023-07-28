# Go-Day6

## 泛型

从1.18版本开始，Go添加了对泛型的支持，也即类型参数。

作为泛型函数的示例，MapKeys接受任意类型的Map并返回其Key的切片。这个函数有2个类型参数 - `K` 和 `V`； `K` 是 `comparable` 类型，也就是说我们可以通过 `==` 和 `!=` 操作符对这个类型的值进行比较。这是Go中Map的Key所必须的。 `V` 是 `any` 类型，意味着它不受任何限制 (`any` 是 `interface{}` 的别名类型).

作为泛型类型的示例，List是一个具有任意类型值的单链表。

可以像在常规类型上一样定义泛型类型的方法，但必须保留类型参数，这个类型是 List[T]，而不是List

当调用泛型函数的时候，经常可以使用类型推断。注意，当调用MapKeys的时候，我们不需要为K和V指定类型-编译器会进行自动推断。

```Go
package main

import "fmt"

func MapKeys[K comparable, V any](m map[K]V) []K {
    r := make([]K, 0, len(m))
    for k := range m {
        r = append(r, k)
    }
    return r
}

type List[T any] struct {
    head, tail *element[T]
}

type element[T any] struct {
    next *element[T]
    val  T
}

func (lst *List[T]) Push(v T) {
    if lst.tail == nil {
        lst.head = &element[T]{val: v}
        lst.tail = lst.head
    } else {
        lst.tail.next = &element[T]{val: v}
        lst.tail = lst.tail.next
    }
}

func (lst *List[T]) GetAll() []T {
    var elems []T
    for e := lst.head; e != nil; e = e.next {
        elems = append(elems, e.val)
    }
    return elems
}

func main() {
    var m = map[int]string{1: "2", 2: "4", 4: "8"}

    fmt.Println("keys m:", MapKeys(m))

    _ = MapKeys[int, string](m)

    lst := List[int]{}
    lst.Push(10)
    lst.Push(13)
    lst.Push(23)
    fmt.Println("list:", lst.GetAll())
}
```

```shell
keys: [4 1 2]
list: [10 13 23]
```

## 错误处理

符合Go语言习惯的做法是使用一个独立、明确的返回值来传递错误信息。这与Java、Ruby使用的异常（exception）以及在C语言中有时用到的重载（overloaded）的单返回/错误值有着明显的不同。

Go语言的处理方式能清楚的知道哪个函数返回了错误，并使用跟其他（无异常处理的）语言类似的方式来处理错误。

按照惯例，错误通常是最后一个返回值并且是error类型，它是一个内建的接口。

errors.New使用给定的错误信息构造一个基本的error值。

返回错误值为nil代表没有错误

还可以通过实现Error()方法来自定义error类型。这里使用自定义错误类型来表示上面例子中的参数错误

```Go
package main

import (
    "errors"
    "fmt"
)

func f1(arg int) (int, error) {
    if arg == 42 {

        return -1, errors.New("can't work with 42")

    }

    return arg + 3, nil
}

type argError struct {
    arg  int
    prob string
}

func (e *argError) Error() string {
    return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

func f2(arg int) (int, error) {
    if arg == 42 {

        return -1, &argError{arg, "can't work with it"}
    }
    return arg + 3, nil
}

func main() {

    for _, i := range []int{7, 42} {
        if r, e := f1(i); e != nil {
            fmt.Println("f1 failed:", e)
        } else {
            fmt.Println("f1 worked:", r)
        }
    }
    for _, i := range []int{7, 42} {
        if r, e := f2(i); e != nil {
            fmt.Println("f2 failed:", e)
        } else {
            fmt.Println("f2 worked:", r)
        }
    }

    _, e := f2(42)
    if ae, ok := e.(*argError); ok {
        fmt.Println(ae.arg)
        fmt.Println(ae.prob)
    }
}
```

在代码中，使用 `&argError`语法来建立一个新的结构体，并提供了arg和prob两个字段的值。

随后的两个循环测试了每一个会返回错误的函数。注意，在if的同一行进行错误检查，是Go代码中的一种常见用法。

```shell
$ go run errors.go
f1 worked: 10
f1 failed: can't work with 42
f2 worked: 10
f2 failed: 42 - can't work with it
42
can't work with it
```

如果想在程序中使用自定义错误类型的数据，需要通过类型断言来得到这个自定义错误类型的实例。

## 协程

协程（goroutine）是轻量级的执行线程。

假设我们有一个函数叫做 `f(s)`。 我们一般会这样 `同步地` 调用它

使用 `go f(s)` 在一个协程中调用这个函数。 这个新的 Go 协程将会 `并发地` 执行这个函数。

也可以为匿名函数启动一个协程

现在两个协程在独立的协程中异步地运行，然后等待两个协程完成（更好的方法是使用WaitGroup）

当运行程序时，首先能看到阻塞式调用的输出，然后是两个协程的交替输出。这种交替的情况表示Go runtime是以并发的方式运行协程的。

```Go
package main

import (
    "fmt"
    "time"
)

func f(from string) {
    for i := 0; i < 3; i++ {
        fmt.Println(from, ":", i)
    }
}

func main() {

    f("direct")

    go f("goroutine")

    go func(msg string) {
        fmt.Println(msg)
    }("going")

    time.Sleep(time.Second)
    fmt.Println("done")
}
```

```shell
$ go run goroutines.go
direct : 0
direct : 1
direct : 2
goroutine : 0
going
goroutine : 1
goroutine : 2
done
```

## 通道

通道（channels）是连接多个协程的管道。可以从一个协程将值发送到通道，然后在另一个协程中接收。

使用make(chan val-type)创建一个新的通道。通道类型就是它们需要传递值的类型

使用channel <- 语法发送一个新的值到通道中。这里我们在一个新的协程中发送“ping“到上面创建的messages通道中。

使用 <-channel 语法从通道中接收一个值，这里我们会收到在上面发送的”ping“消息并将其打印出来。

运行程序时，通过通道，成功的将消息”ping“从一个协程传送到了另一个协程中。

默认发送和接收操作是阻塞的，直到发送方和接收方都就绪，这个特性允许我们，不适用任何其它的同步操作，就可以在程序结尾处等待消息”ping“。

```Go
package main

import "fmt"

func main() {

    messages := make(chan string)

    go func() { messages <- "ping" }()

    msg := <-messages
    fmt.Println(msg)
}
```

```shell
$ go run channels.go
ping
```

