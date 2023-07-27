# Go-Day12

## Recover

Go通过使用recover内置函数，可以从panic中恢复recover。recover可以阻止panic中止程序，并让它继续执行。

在这样的例子中很有用：当其中一个客户端连接出现严重错误，服务器不希望崩溃。相反，服务器希望关闭该连接并继续为其他的客户端提供服务。实际上，这就是Go的net/http包默认对于HTTP服务器的处理。

```Go
package main

import "fmt"

func mayPanic() {
    panic("a problem")
}

func main() {

    defer func() {
        if r := recover(); r != nil {

            fmt.Println("Recovered. Error:\n", r)
        }
    }()

    mayPanic()

    fmt.Println("After mayPanic()")
}
```

代码中写了一个panic函数

必须在defer函数中调用recover。当跳出引发panic的函数时，defer会被激活，其中的recover会捕获panic。

recover的返回值是在调用panic时抛出的错误。

最后一行代码不会执行，因为mayPanic函数会调用panic。main程序的执行在panic点停止，并在继续处理完defer后结束。

```shell
Recovered. Error:
 a problem
```

## 字符串函数

标准库的 `strings` 包提供了很多有用的字符串相关的函数。 这儿有一些用来让你对 `strings` 包有一个初步了解的例子。

```Go
package main

import (
    "fmt"
    s "strings"
)

var p = fmt.Println

func main() {

    p("Contains:  ", s.Contains("test", "es"))
    p("Count:     ", s.Count("test", "t"))
    p("HasPrefix: ", s.HasPrefix("test", "te"))
    p("HasSuffix: ", s.HasSuffix("test", "st"))
    p("Index:     ", s.Index("test", "e"))
    p("Join:      ", s.Join([]string{"a", "b"}, "-"))
    p("Repeat:    ", s.Repeat("a", 5))
    p("Replace:   ", s.Replace("foo", "o", "0", -1))
    p("Replace:   ", s.Replace("foo", "o", "0", 1))
    p("Split:     ", s.Split("a-b-c-d-e", "-"))
    p("ToLower:   ", s.ToLower("TEST"))
    p("ToUpper:   ", s.ToUpper("test"))
    p()

    p("Len: ", len("hello"))
    p("Char:", "hello"[1])
}
```

给 `fmt.Println` 一个较短的别名， 因为我们随后会大量的使用它。

这是一些 `strings` 中有用的函数例子。 由于它们都是包的函数，而不是字符串对象自身的方法， 这意味着我们需要在调用函数时，将字符串作为第一个参数进行传递。 你可以在 [`strings`](http://golang.org/pkg/strings/) 包文档中找到更多的函数。

虽然不是 `strings` 的函数，但仍然值得一提的是， 获取字符串长度（以字节为单位）以及通过索引获取一个字节的机制。

注意，上面的 `len` 以及索引工作在字节级别上。 Go 使用 UTF-8 编码字符串，因此通常按原样使用。 如果您可能使用多字节的字符，则需要使用可识别编码的操作。 详情请参考 [strings, bytes, runes and characters in Go](https://blog.golang.org/strings)。

```shell
	
$ go run string-functions.go
Contains:   true
Count:      2
HasPrefix:  true
HasSuffix:  true
Index:      1
Join:       a-b
Repeat:     aaaaa
Replace:    f00
Replace:    f0o
Split:      [a b c d e]
ToLower:    test
ToUpper:    TEST
Len:  5
Char: 101
```

## 字符串格式化

Go 在传统的 `printf` 中对字符串格式化提供了优异的支持。 这儿有一些基本的字符串格式化的任务的例子。

Go 提供了一些用于格式化常规值的打印“动词”。 例如，这样打印 `point` 结构体的实例。

如果值是一个结构体，`%+v` 的格式化输出内容将包括结构体的字段名。

`%#v` 根据 Go 语法输出值，即会产生该值的源码片段。

需要打印值的类型，使用 `%T`。

格式化布尔值很简单。

格式化整型数有多种方式，使用 `%d` 进行标准的十进制格式化。

这个输出二进制表示形式。

输出给定整数的对应字符。

`%x` 提供了十六进制编码。

同样的，也为浮点型提供了多种格式化选项。 使用 `%f` 进行最基本的十进制格式化。

`%e` 和 `%E` 将浮点型格式化为（稍微有一点不同的）科学记数法表示形式。

使用 `%s` 进行基本的字符串输出。

像 Go 源代码中那样带有双引号的输出，使用 `%q`。

和上面的整型数一样，`%x` 输出使用 base-16 编码的字符串， 每个字节使用 2 个字符表示。

要输出一个指针的值，使用 `%p`。

格式化数字时，您经常会希望控制输出结果的宽度和精度。 要指定整数的宽度，请在动词 “%” 之后使用数字。 默认情况下，结果会右对齐并用空格填充。

你也可以指定浮点型的输出宽度，同时也可以通过 `宽度.精度` 的语法来指定输出的精度。

要左对齐，使用 `-` 标志。

也许也想控制字符串输出时的宽度，特别是要确保他们在类表格输出时的对齐。 这是基本的宽度右对齐方法。

要左对齐，和数字一样，使用 `-` 标志。

到目前为止，我们已经看过 `Printf` 了， 它通过 `os.Stdout` 输出格式化的字符串。 `Sprintf` 则格式化并返回一个字符串而没有任何输出。

你可以使用 `Fprintf` 来格式化并输出到 `io.Writers` 而不是 `os.Stdout`。

```Go
package main

import (
    "fmt"
    "os"
)

type point struct {
    x, y int
}

func main() {

    p := point{1, 2}
    fmt.Printf("struct1: %v\n", p)

    fmt.Printf("struct2: %+v\n", p)

    fmt.Printf("struct3: %#v\n", p)

    fmt.Printf("type: %T\n", p)

    fmt.Printf("bool: %t\n", true)

    fmt.Printf("int: %d\n", 123)

    fmt.Printf("bin: %b\n", 14)

    fmt.Printf("char: %c\n", 33)

    fmt.Printf("hex: %x\n", 456)

    fmt.Printf("float1: %f\n", 78.9)

    fmt.Printf("float2: %e\n", 123400000.0)
    fmt.Printf("float3: %E\n", 123400000.0)

    fmt.Printf("str1: %s\n", "\"string\"")

    fmt.Printf("str2: %q\n", "\"string\"")

    fmt.Printf("str3: %x\n", "hex this")

    fmt.Printf("pointer: %p\n", &p)

    fmt.Printf("width1: |%6d|%6d|\n", 12, 345)

    fmt.Printf("width2: |%6.2f|%6.2f|\n", 1.2, 3.45)

    fmt.Printf("width3: |%-6.2f|%-6.2f|\n", 1.2, 3.45)

    fmt.Printf("width4: |%6s|%6s|\n", "foo", "b")

    fmt.Printf("width5: |%-6s|%-6s|\n", "foo", "b")

    s := fmt.Sprintf("sprintf: a %s", "string")
    fmt.Println(s)

    fmt.Fprintf(os.Stderr, "io: an %s\n", "error")
}
```

```shell
$ go run string-formatting.go
struct1: {1 2}
struct2: {x:1 y:2}
struct3: main.point{x:1, y:2}
type: main.point
bool: true
int: 123
bin: 1110
char: !
hex: 1c8
float1: 78.900000
float2: 1.234000e+08
float3: 1.234000E+08
str1: "string"
str2: "\"string\""
str3: 6865782074686973
pointer: 0xc0000ba000
width1: |    12|   345|
width2: |  1.20|  3.45|
width3: |1.20  |3.45  |
width4: |   foo|     b|
width5: |foo   |b     |
sprintf: a string
io: an error
```

