# Go-Day2

## Switch

Switch是多分支情况时快捷的条件语句

一个基本的Switch语句程序如下：

```Go
package main

import (
    "fmt"
    "time"
)

func main() {

    i := 2
    fmt.Print("write ", i, " as ")
    switch i {
    case 1:
        fmt.Println("one")
    case 2:
        fmt.Println("two")
    case 3:
        fmt.Println("three")
    }

    switch time.Now().Weekday() {
    case time.Saturday, time.Sunday:
        fmt.Println("It's the weekend")
    default:
        fmt.Println("It's a weekday")
    }

    t := time.Now()
    switch {
    case t.Hour() < 12:
        fmt.Println("It's before noon")
    default:
        fmt.Println("It's after noon")
    }

    whatAmI := func(i interface{}) {
        switch t := i.(type) {
        case bool:
            fmt.Println("I'm a bool")
        case int:
            fmt.Println("I'm an int")
        default:
            fmt.Printf("Don't know type %T\n", t)
        }
    }
    whatAmI(true)
    whatAmI(1)
    whatAmI("hey")
}
```

```shell
$ go run switch.go
Write 2 as two
It's a weekday
It's after noon
I'm a bool
I'm an int
Don't know type string
```

在同一个case语句中，可以使用逗号来分隔多个表达式，也可以使用可选的default分支。

不带表达式的switch是实现if/else逻辑的另一种方式。case表达式也可以不使用常量

类型开关（type switch）比较类型而非值。可以用来发现一个接口值的类型。

## 数组

在Go中，数组是一个具有编号且长度固定的元素序列

```Go
package main

import "fmt"

func main() {

    var a [5]int
    fmt.Println("emp:", a)

    a[4] = 100
    fmt.Println("set:", a)
    fmt.Println("get:", a[4])

    fmt.Println("len:", len(a))

    b := [5]int{1, 2, 3, 4, 5}
    fmt.Println("dcl:", b)

    var twoD [2][3]int
    for i := 0; i < 2; i++ {
        for j := 0; j < 3; j++ {
            twoD[i][j] = i + j
        }
    }
    fmt.Println("2d: ", twoD)
}
```

在例子中，创建了一个刚好可以存放5个int元素的数组a。元素的类型和长度都是数组类型的一部分。数组默认值是零值，对于int数组来说，元素的零值是0。

可以使用 `array[index]=value`语法来设置数组指定位置的值，或者用 `array[index]`来得到值

内置函数 `len`可以返回数组的长度

数组类型是一维的，但是可以组个构造多维的数据结构。

使用 `fmt.Println()`打印数组时，会按照 `v1 v2 v3...`的格式打印

```shell
$ go run arrays.go
emp: [0 0 0 0 0]
set: [0 0 0 0 100]
get: 100
len: 5
dcl: [1 2 3 4 5]
2d:  [[0 1 2] [1 2 3]]
```

## 切片

Slice是Go中一个重要的数据类型，它提供了比数组更强大的序列交互方式。

与数组不同，Slice的类型仅由它所包含的元素的类型决定（与元素个数无关）。要创建一个长度不为0的空Slice，需要使用内建函数make。

和数组一样的方法

除了基本操作外，slice支持比数组更丰富的操作，比如slice支持内建函数append，该函数会返回一个包含了一个或者多个新值的slice。注意由于append可能返回一个新的slice，我们需要接收其返回值。

slice还可以copy。

slice支持通过 `slice[low:high]`语法进行“切片”操作。

也可以在一行代码中声明并初始化一个slice变量

Slice也可以组成多维数据结构，内部的slice长度可以不一致，这一点和多维数组不同。

```go
package main

import "fmt"

func main() {

    s := make([]string, 3)
    fmt.Println("emp:", s)

    s[0] = "a"
    s[1] = "b"
    s[2] = "c"
    fmt.Println("set:", s)
    fmt.Println("get:", s[2])

    fmt.Println("len:", len(s))

    s = append(s, "d")
    s = append(s, "e", "f")
    fmt.Println("apd:", s)

    c := make([]string, len(s))
    copy(c, s)
    fmt.Println("cpy:", c)

    l := s[2:5]
    fmt.Println("sl1:", l)

    l = s[:5]
    fmt.Println("sl2:", l)

    l = s[2:]
    fmt.Println("sl3:", l)

    t := []string{"g", "h", "i"}
    fmt.Println("dcl:", t)

    twoD := make([][]int, 3)
    for i := 0; i < 3; i++ {
        innerLen := i + 1
        twoD[i] = make([]int, innerLen)
        for j := 0; j < innerLen; j++ {
            twoD[i][j] = i + j
        }
    }
    fmt.Println("2d: ", twoD)
}
```

Slice和数组是不同的类型，但他们通过 `fmt.Println`打印的输出结果是类似的

```shell
$ go run slices.go
emp: [  ]
set: [a b c]
get: c
len: 3
apd: [a b c d e f]
cpy: [a b c d e f]
sl1: [c d e]
sl2: [a b c d e]
sl3: [c d e f]
dcl: [g h i]
2d:  [[0] [1 2] [2 3 4]]
```



## Map

Map是Go内建的关联数据类型（又称哈希hash或者字典dict）

要创建一个空map，需要使用内建函数make：`make(map[key-type]val-type)`

使用典型的 `name[key]=val`语法来设置键值对

打印map，例如，使用fmt.Println打印一个map，会输出它所有的键值对。

使用 `name[key]`来获取一个键的值

内建函数 `len`可以返回一个map的键值对数量

内建函数 `delete`可以从一个map中移除键值对

```Go
package main

import "fmt"

func main() {

    m := make(map[string]int)

    m["k1"] = 7
    m["k2"] = 13

    fmt.Println("map:", m)

    v1 := m["k1"]
    fmt.Println("v1: ", v1)

    fmt.Println("len:", len(m))

    delete(m, "k2")
    fmt.Println("map:", m)

    _, prs := m["k2"]
    fmt.Println("prs:", prs)

    n := map[string]int{"foo": 1, "bar": 2}
    fmt.Println("map:", n)
}
```

当从一个map中取值时，还有可以选择是否接收的第二个返回值，该值表明了map中是否存在这个键。这可以用来消除 `键不存在`和 `键的值为零值`所产生的歧义，例如0和“”。这里不需要值的话，可以用 `空白标识符 _`将其忽略

在使用 `fmt.Println`方法打印一个map时，是以 `map[k:v k:v]`的格式输出的

```shell
$ go run maps.go 
map: map[k1:7 k2:13]
v1:  7
len: 2
map: map[k1:7]
prs: false
map: map[foo:1 bar:2]
```

