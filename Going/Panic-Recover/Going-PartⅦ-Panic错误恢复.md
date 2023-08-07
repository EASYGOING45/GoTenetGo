# Going-PartⅦ-Panic错误恢复

- 实现错误处理机制
- 框架收尾

## Panic

Go语言中，比较常见的错误处理方法是返回error，由调用者决定后续如何处理。但是如果是无法恢复的错误，可以手动触发panic，当然如果在程序运行过程中出现了类似于数组越界的错误，panic也会被触发。panic会中止当前执行的程序，退出：

主动触发例子如下：

```Go
// hello.go
func main() {
	fmt.Println("before panic")
	panic("crash")
	fmt.Println("after panic")
}
```

```shell
$ go run hello.go

before panic
panic: crash

goroutine 1 [running]:
main.main()
        ~/go_demo/hello/hello.go:7 +0x95
exit status 2
```

数组越界触发的panic

```Go
// hello.go
func main() {
	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```

```shell
go run hello.go
panic: runtime error: index out of range [4] with length 3
```

## Defer

Panic会导致程序被中止，但是在退出前，会先处理完当前协程上已经Defer的任务，执行完成后再退出。效果类似于java语言的 `try...catch`

```Go
// hello.go
func main() {
	defer func() {
		fmt.Println("defer func")
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```

```shell
$ go run hello.go 
defer func
panic: runtime error: index out of range [4] with length 3
```

可以defer多个任务，在同一个函数中defer多个任务，会逆序执行，即先执行最后defer的任务。

在这里，defer的任务执行完成之后，panic还会继续被抛出，导致程序非正常结束	

## Revocer

Go语言还提供了recover函数，可以避免因为Panic发生而导致整个程序终止，recover函数只在defer中生效。

```Go
// hello.go
func test_recover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}

func main() {
	test_recover()
	fmt.Println("after recover")
}
```

```shell
$ go run hello.go 
defer func
recover success
after recover
```

可以看到，recover捕获了panic，程序正常结束。`test_recover()`中的after panic没有打印，这是正确的，当panic被触发时，控制权就被交给了defer，就像在java中，try代码块中发生了异常，控制权交给了`catch`，接下来执行`catch`代码块中的代码。而在main()中打印了after recover，说明程序已经恢复正常，继续往下执行直到结束。

## Going的错误处理机制

对一个Web框架而言，错误处理机制是非常必要的，可能是框架本身没有完备的测试，导致在某些情况下出现空指针异常等情况。也有可能用户不正确的参数，触发了某些异常，例如数组越界，空指针等。如果因为这些原因导致系统宕机，必然是不可接受的。

此前框架中没有加入异常处理机制，如果代码中存在会触发panic的bug，很容易宕掉。

如下：

```Go
func main() {
	r := gee.New()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
```

在上面的代码中，我们为 gee 注册了路由 `/panic`，而这个路由的处理函数内部存在数组越界 `names[100]`，如果访问 *localhost:9999/panic*，Web 服务就会宕掉。

今天，我们将在 gee 中添加一个非常简单的错误处理机制，即在此类错误发生时，向用户返回 *Internal Server Error*，并且在日志中打印必要的错误信息，方便进行错误定位。

我们之前实现了中间件机制，错误处理也可以作为一个中间件，增强 going 框架的能力。

新增文件 **going/recovery.go**，在这个文件中实现中间件 `Recovery`。

```Go
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
```

`Recovery` 的实现非常简单，使用 defer 挂载上错误恢复的函数，在这个函数中调用 *recover()*，捕获 panic，并且将堆栈信息打印在日志中，向用户返回 *Internal Server Error*。

你可能注意到，这里有一个 *trace()* 函数，这个函数是用来获取触发 panic 的堆栈信息，完整代码如下：

```Go
package going

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
```

在 *trace()* 中，调用了 `runtime.Callers(3, pcs[:])`，Callers 用来返回调用栈的程序计数器, 第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 `defer func`。因此，为了日志简洁一点，我们跳过了前 3 个 Caller。

接下来，通过 `runtime.FuncForPC(pc)` 获取对应的函数，在通过 `fn.FileLine(pc)` 获取到调用该函数的文件名和行号，打印在日志中。

至此，gee 框架的错误处理机制就完成了。

```Go
package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *going.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
```

接下来进行测试，先访问主页，访问一个有BUG的 `/panic`，服务正常返回。接下来我们再一次成功访问了主页，说明服务完全运转正常。

```shell
$ curl "http://localhost:9999"
Hello Geektutu
$ curl "http://localhost:9999/panic"
{"message":"Internal Server Error"}
$ curl "http://localhost:9999"
Hello Geektutu
```

我们可以在后台日志中看到如下内容，引发错误的原因和堆栈信息都被打印了出来，通过日志，我们可以很容易地知道，在*day7-panic-recover/main.go:47* 的地方出现了 `index out of range` 错误。

![image-20230807162629916](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230807162629916.png)