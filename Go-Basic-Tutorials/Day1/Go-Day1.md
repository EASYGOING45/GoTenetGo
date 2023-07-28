# Go-Day1

## Hello World

第一个程序将打印传说中的“hello world”

要运行这个程序，先将将代码放到名为 `hello-world.go` 的文件中，然后执行 `go run`。

如果我们想将程序编译成二进制文件（Windows 平台是 .exe 可执行文件）， 可以通过 `go build` 来达到目的。

然后我们可以直接运行这个二进制文件。

```Go
package main

import "fmt"

func main() {
    fmt.Println("hello world")
}
```

```shell
$ go run hello-world.go
hello world

$ go build hello-world.go
$ ls
hello-world    hello-world.go

$ ./hello-world
hello world
```

## 值

GO拥有多种值类型，包括字符串、整型、浮点型、布尔型等。

```Go
package main

import "fmt"

func main() {

    fmt.Println("go" + "lang")

    fmt.Println("1+1 =", 1+1)
    fmt.Println("7.0/3.0 =", 7.0/3.0)

    fmt.Println(true && false)
    fmt.Println(true || false)
    fmt.Println(!true)
}
```

字符串可以通过+连接。

整数和浮点数，布尔型，以及常见的布尔操作

```shell
$ go run values.go
golang
1+1 = 2
7.0/3.0 = 2.3333333333333335
false
true
false
```

## 变量

在Go中，变量需要显式声明，并且在函数调用等情况下，编译器会检查其类型的正确性。

```Go
package main

import "fmt"

func main() {

    var a = "initial"
    fmt.Println(a)

    var b, c int = 1, 2
    fmt.Println(b, c)

    var d = true
    fmt.Println(d)

    var e int
    fmt.Println(e)

    f := "short"
    fmt.Println(f)
}
```

var声明1个或者多个变量

可以一次性声明多个变量

Go会自动推断已经有初始值的变量的类型

声明后却没有给出对应的初始值时，变量将会初始化为零值。例如，int的零值是0.

> :=语法是声明并初始化变量的简写，例如 `var f string = short`可以简写为 `f := "short"`

```shell
$ go run variables.go
initial
1 2
true
0
short
```

## 常量

Go支持字符、字符串、布尔和数值常量。

const用于声明一个常量。

```Go
package main

import (
    "fmt"
    "math"
)

const s string = "constant"

func main() {
    fmt.Println(s)

    const n = 500000000

    const d = 3e20 / n
    fmt.Println(d)

    fmt.Println(int64(d))

    fmt.Println(math.Sin(n))
}
```

const语句可以出现在任何var语句可以出现的地方

常数表达式可以执行任意精度的运算

数值型常量没有确定的类型，直到被给定某个类型，比如显式类型转化。

一个数字可以根据上下文的需要（比如变量赋值、函数调用）自动确定类型。举例：math.Sin函数需要一个float64的参数，n会自动确定类型。

```shell
$ go run constant.go 
constant
6e+11
600000000000
-0.28470407323754404
```

## For循环

for是Go中唯一的循环结构。展示一些for的基本使用方式

最基础的方式：单个循环条件

经典的初始/条件/后续for循环

不带条件的for循环将一直重复执行，直到在循环体内使用了break或者return跳出循环

也可以使用continue直接进入下一次循环

```Go
package main

import "fmt"

func main() {

    i := 1
    for i <= 3 {
        fmt.Println(i)
        i = i + 1
    }

    for j := 7; j <= 9; j++ {
        fmt.Println(j)
    }

    for {
        fmt.Println("loop")
        break
    }

    for n := 0; n <= 5; n++ {
        if n%2 == 0 {
            continue
        }
        fmt.Println(n)
    }
}
```

后续会接触到range语句，channels以及其他数据结构时，会有一些for的其他用法

```shell
$ go run for.go
1
2
3
7
8
9
loop
1
3
5
```

