# Go-Day3

## Range遍历

range用于迭代各种各样的数据结构。可以使用range来对slice中的元素求和。数组也可以用这种方法初始化并赋初值。

range在数组和slice中提供对每项的索引和值的访问，也可以使用空白标识符 _将其忽略。实际上，有时是需要这个索引的

```Go
package main

import "fmt"

func main() {

    nums := []int{2, 3, 4}
    sum := 0
    for _, num := range nums {
        sum += num
    }
    fmt.Println("sum:", sum)

    for i, num := range nums {
        if num == 3 {
            fmt.Println("index:", i)
        }
    }

    kvs := map[string]string{"a": "apple", "b": "banana"}
    for k, v := range kvs {
        fmt.Printf("%s -> %s\n", k, v)
    }

    for k := range kvs {
        fmt.Println("key:", k)
    }

    for i, c := range "go" {
        fmt.Println(i, c)
    }
}
```

range可以在map中迭代键值对，也可以只遍历map的键，range在字符串中迭代unicode码点(code point)，第一个返回值是字符的起始字节位置，第二个是字符本身。

```shell
$ go run range.go
sum: 9
index: 1
a -> apple
b -> banana
key: a
key: b
0 103
1 111
```

## 函数

函数是Go的核心

在以下代码中，写了一个函数，接收两个int并且以int返回它们的和

Go需要明确的return，也就是说，它不会自动return最后一个表达式的值

当多个连续的参数为同样类型时，可以仅声明最后一个参数的类型，忽略之前相同类型参数的类型声明

也可以通过函数名（参数列表）的方法来调用函数

Go中的函数还有很多其他的特性，期中一个就是多值返回。

```Go
package main

import "fmt"

func plus(a int, b int) int {

    return a + b
}

func plusPlus(a, b, c int) int {
    return a + b + c
}

func main() {

    res := plus(1, 2)
    fmt.Println("1+2 =", res)

    res = plusPlus(1, 2, 3)
    fmt.Println("1+2+3 =", res)
}
```

```shell
$ go run functions.go
1+2 = 3
1+2+3 = 6
```

## 多返回值

Go原生支持多返回值，这个特性在Go中经常用到，例如用来同时返回一个函数的结果和错误信息。

(int,int)在函数中标志着这个函数返回2个int

可以通过多赋值操作来使用这两个不同的返回值

如果仅仅需要返回值的一部分的话，可以使用空白标识符 _ 。

```Go
package main

import "fmt"

func vals() (int, int) {
    return 3, 7
}

func main() {

    a, b := vals()
    fmt.Println(a)
    fmt.Println(b)

    _, c := vals()
    fmt.Println(c)
}
```

```shell
$ go run multiple-return-values.go
3
7
7
```

## 变参函数

又称可变参数函数，在调用时可以传递任意数量的参数，例如 fmt.Println就是一个常见的变参函数

这个函数接受任意数量的int作为参数

变参函数使用常规的调用方式，传入独立的参数

如果有一个含有多个值的slice，想把它们作为参数使用，需要这样调用 func(slice...)

```Go
package main

import "fmt"

func sum(nums ...int) {
    fmt.Print(nums, " ")
    total := 0
    for _, num := range nums {
        total += num
    }
    fmt.Println(total)
}

func main() {

    sum(1, 2)
    sum(1, 2, 3)

    nums := []int{1, 2, 3, 4}
    sum(nums...)
}
```

```shell
$ go run variadic-functions.go 
[1 2] 3
[1 2 3] 6
[1 2 3 4] 10
```

