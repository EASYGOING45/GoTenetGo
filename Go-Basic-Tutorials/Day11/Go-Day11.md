# Go-Day11

## 排序

Go的sort包实现了内建及用户自定义数据类型的排序功能。

```Go
package main

import (
    "fmt"
    "sort"
)

func main() {

    strs := []string{"c", "a", "b"}
    sort.Strings(strs)
    fmt.Println("Strings:", strs)

    ints := []int{7, 2, 4}
    sort.Ints(ints)
    fmt.Println("Ints:   ", ints)

    s := sort.IntsAreSorted(ints)
    fmt.Println("Sorted: ", s)
}
```

排序方法是针对内置数据类型的，这是一个字符串排序的例子。注意，这是原地排序的，所以会直接改变给定的切片，而非返回一个新切片。

代码展示了一个int排序的例子

也可以使用sort来检查一个切片是否为有序的。

运行程序，打印排序好的字符串和整型切片，以及数组是否有序的检查结果：true。

```shell
$ go run sorting.go
Strings: [a b c]
Ints:    [2 4 7]
Sorted:  true
```

## 使用函数自定义排序

有时，可能想根据自然顺序以外的方式来对集合进行排序。 例如，假设我们要按字符串的长度而不是按字母顺序对它们进行排序。 这儿有一个在 Go 中自定义排序的示例。

为了在 Go 中使用自定义函数进行排序，我们需要一个对应的类型。 我们在这里创建了一个 `byLength` 类型，它只是内建类型 `[]string` 的别名。

```Go
package main

import (
    "fmt"
    "sort"
)

type byLength []string

func (s byLength) Len() int {
    return len(s)
}
func (s byLength) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
    return len(s[i]) < len(s[j])
}

func main() {
    fruits := []string{"peach", "banana", "kiwi"}
    sort.Sort(byLength(fruits))
    fmt.Println(fruits)
}
```

我们为该类型实现了 `sort.Interface` 接口的 `Len`、`Less` 和 `Swap` 方法， 这样我们就可以使用 `sort` 包的通用 `Sort` 方法了， `Len` 和 `Swap` 在各个类型中的实现都差不多， `Less` 将控制实际的自定义排序逻辑。 在这个的例子中，我们想按字符串长度递增的顺序来排序， 所以这里使用了 `len(s[i])` 和 `len(s[j])` 来实现 `Less`。

一切准备就绪后，我们就可以通过将切片 `fruits` 强转为 `byLength` 类型的切片， 然后对该切片使用 `sort.Sort` 来实现自定义排序。

运行这个程序，和预期的一样， 显示了一个按照字符串长度排序的列表。

类似的，参照这个例子，创建一个自定义类型， 为它实现 `Interface` 接口的三个方法， 然后对这个自定义类型的切片调用 `sort.Sort` 方法， 我们就可以通过任意函数对 Go 切片进行排序了。

```shell
$ go run sorting-by-functions.go 
[kiwi peach banana]
```

## Panic

panic意味着有些出乎意料的错误发生，通常用它来表示程序正常运行中不应该出现的错误，或者我们不准备优雅处理的错误。

代码中使用panic来检查这个站点上预期之外的错误，而该站点上只有一个程序：触发panic。

panic的一种常见用法是：当函数返回我们不知道如何处理（或不象处理）的错误值时，中止操作。如果创建新文件时遇到意外错误该如何处理？代码展示了panic的示例。

```Go
package main

import "os"

func main() {

    panic("a problem")

    _, err := os.Create("/tmp/file")
    if err != nil {
        panic(err)
    }
}
```

运行程序将会导致panic：输出一个错误消息和协程追踪信息，并以非零的状态退出程序。

当main中触发第一个panic时，程序就会退出而不会执行代码的其余部分，如果想看到程序尝试创建/tmp/file文件，请注释掉第一个panic。

```shell
$ go run panic.go
panic: a problem
goroutine 1 [running]:
main.main()
    /.../panic.go:12 +0x47
...
exit status 2
```

注意，与某些使用exception处理错误的语言不同，在Go中，通常会尽可能的使用返回值来标示错误。

## Defer

Defer用于确保程序在执行完成后，会调用某个函数，一般是执行清理工作，Defer的用途跟其他语言的ensure或finally类似。

假设我们想要创建一个文件、然后写入数据、最后在程序结束时关闭该文件。代码展示了如何通过defer来做到这一切。

在createFile后立即得到一个文件对象，我们使用defer通过closeFile来关闭这个文件。这会在封闭函数(main)结束时执行，即writeFile完成以后。

```Go
package main

import (
    "fmt"
    "os"
)

func main() {

    f := createFile("/tmp/defer.txt")
    defer closeFile(f)
    writeFile(f)
}

func createFile(p string) *os.File {
    fmt.Println("creating")
    f, err := os.Create(p)
    if err != nil {
        panic(err)
    }
    return f
}

func writeFile(f *os.File) {
    fmt.Println("writing")
    fmt.Fprintln(f, "data")

}

func closeFile(f *os.File) {
    fmt.Println("closing")
    err := f.Close()

    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
```

关闭文件时，进行错误检查是非常重要的，即使在defer函数中也是如此。

执行程序，确认写入数据后，文件已关闭。

```shell
go run defer.go
creating
writing
closing
```

