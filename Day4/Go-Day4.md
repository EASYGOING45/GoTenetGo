# Go-Day4

## 闭包

Go支持匿名函数，并能用其构造闭包，匿名函数在想定义一个不需要命名的内联函数时是很使用的

以下定义的intSeq函数返回一个在其函数体内定义的匿名函数。返回的函数使用闭包的方式隐藏变量i，返回的函数隐藏变量i以形成闭包

我们调用intSeq函数，将返回值（一个函数）赋给nextInt,这个函数的值包含了自己的值i，这样在每次调用nextInt时，都会更新i的值。

通过多次调用来查看闭包的效果。

```Go
package main

import "fmt"

func intSeq() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

func main() {

    nextInt := intSeq()

    fmt.Println(nextInt())
    fmt.Println(nextInt())
    fmt.Println(nextInt())

    newInts := intSeq()
    fmt.Println(newInts())
}
```

```shell
$ go run closures.go
1
2
3
1
```

## 递归

Go支持递归

fact函数在到达fact(0)前一直调用自身

闭包也可以是递归的，但这要求在定义闭包之前用类型化的var显式声明闭包。

由于fib之前在main中声明过，因此Go知道在这里用fib调用哪个函数。

```Go
package main

import "fmt"

func fact(n int) int {
    if n == 0 {
        return 1
    }
    return n * fact(n-1)
}

func main() {
    fmt.Println(fact(7))

    var fib func(n int) int

    fib = func(n int) int {
        if n < 2 {
            return n
        }
        return fib(n-1) + fib(n-2)

    }

    fmt.Println(fib(7))
}
```

```shell
$ go run recursion.go 
5040
13
```

## 指针

Go支持指针，允许在程序中通过引用传递来传递值和数据结构

在编写的代码中，通过两个函数：`zeroval` 和 `zeroptr` 来比较 `指针` 和 `值`。 `zeroval` 有一个 `int` 型参数，所以使用值传递。 `zeroval` 将从调用它的那个函数中得到一个实参的拷贝：ival。

`zeroptr` 有一个和上面不同的参数：`*int`，这意味着它使用了 `int` 指针。 紧接着，函数体内的 `*iptr` 会 *解引用* 这个指针，从它的内存地址得到这个地址当前对应的值。 对解引用的指针赋值，会改变这个指针引用的真实地址的值。

通过 `&i` 语法来取得 `i` 的内存地址，即指向 `i` 的指针。

指针也是可以被打印的。

`zeroval` 在 `main` 函数中不能改变 `i` 的值， 但是 `zeroptr` 可以，因为它有这个变量的内存地址的引用。

```Go
package main

import "fmt"

func zeroval(ival int) {
    ival = 0
}

func zeroptr(iptr *int) {
    *iptr = 0
}

func main() {
    i := 1
    fmt.Println("initial:", i)

    zeroval(i)
    fmt.Println("zeroval:", i)

    zeroptr(&i)
    fmt.Println("zeroptr:", i)

    fmt.Println("pointer:", &i)
}
```

```shell
$ go run pointers.go
initial: 1
zeroval: 1
zeroptr: 0
pointer: 0x42131100
```

## 字符串和rune类型

Go语言中的字符串是一个只读的byte类型的切片。Go语言和标准库特别对待字符串-作为以UTF-8为编码的文本容器。在其他语言当中，字符串由“字符”组成。在Go语言当中，字符的概念被称为rune-它是一个表示Unicode编码的整数。

在下述代码中，s是一个string，分配了一个literal value表示泰语中的单词“hello”。Go字符串是UTF-8编码的文本

因为字符串等价于[]byte，这会产生存储在其中的原始字节的长度。

对字符串进行索引会在每个索引处生成原始字节值。这个循环生成构成s中Unicode的所有字节的十六进制值。

要计算字符串中有多少rune，我们可以使用utf8包。注意，RuneCountIntString的运行时取决于字符串的大小。因为它必须按顺序解码每个UTF-8 rune。一些泰语字符由多个UTF-8 code point表示。

range循环专门处理字符串并解码每个rune及其在字符串中的偏移量

可以通过显式使用utf8.DecodeRuneInString函数来实现相同的迭代

用单引号括起来的值是rune literals，可以直接将rune value与rune literal进行比较

```Go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {

    const s = "สวัสดี"

    fmt.Println("Len:", len(s))

    for i := 0; i < len(s); i++ {
        fmt.Printf("%x ", s[i])
    }
    fmt.Println()

    fmt.Println("Rune count:", utf8.RuneCountInString(s))

    for idx, runeValue := range s {
        fmt.Printf("%#U starts at %d\n", runeValue, idx)
    }

    fmt.Println("\nUsing DecodeRuneInString")
    for i, w := 0, 0; i < len(s); i += w {
        runeValue, width := utf8.DecodeRuneInString(s[i:])
        fmt.Printf("%#U starts at %d\n", runeValue, i)
        w = width

        examineRune(runeValue)
    }
}

func examineRune(r rune) {

    if r == 't' {
        fmt.Println("found tee")
    } else if r == 'ส' {
        fmt.Println("found so sua")
    }
}
```

```shell
$ go run strings-and-runes.go
Len: 18
e0 b8 aa e0 b8 a7 e0 b8 b1 e0 b8 aa e0 b8 94 e0 b8 b5
Rune count: 6
U+0E2A 'ส' starts at 0
U+0E27 'ว' starts at 3
U+0E31 'ั' starts at 6
U+0E2A 'ส' starts at 9
U+0E14 'ด' starts at 12
U+0E35 'ี' starts at 15
Using DecodeRuneInString
U+0E2A 'ส' starts at 0
found so sua
U+0E27 'ว' starts at 3
U+0E31 'ั' starts at 6
U+0E2A 'ส' starts at 9
found so sua
U+0E14 'ด' starts at 12
U+0E35 'ี' starts at 15
```

