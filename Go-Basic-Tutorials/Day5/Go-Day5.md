# Go-Day5

## 结构体

Go的结构体（struct）是带类型的字段（fields）集合。这在组织数据时非常有用。

代码中的person结构体包含了name和age两个字段。

newPerson使用给定的名字构造一个新的person结构体。

可以安全地返回指向局部变量的指针，因为局部变量将在函数的作用域中继续存在。

使用person{}语法创建新的结构体元素

可以在初始化一个结构体元素时指定字段名字

省略的字段将被初始化为零值

&前缀生成一个结构体指针

在构造函数中封装创建新的结构实例是一种习惯用法

使用.来访问结构体字段，就像C++那样，也可以对结构体指针使用.-指针会被自动解引用。

结构体是可变的（mutable）。

```Go
package main

import "fmt"

type person struct {
    name string
    age  int
}

func newPerson(name string) *person {

    p := person{name: name}
    p.age = 42
    return &p
}

func main() {

    fmt.Println(person{"Bob", 20})

    fmt.Println(person{name: "Alice", age: 30})

    fmt.Println(person{name: "Fred"})

    fmt.Println(&person{name: "Ann", age: 40})

    fmt.Println(newPerson("Jon"))

    s := person{name: "Sean", age: 50}
    fmt.Println(s.name)

    sp := &s
    fmt.Println(sp.age)

    sp.age = 51
    fmt.Println(sp.age)
}
```

```shell
$ go run structs.go
{Bob 20}
{Alice 30}
{Fred 0}
&{Ann 40}
&{Jon 42}
Sean
50
51
```

## 方法

Go支持为结构体类型定义方法（methods）

这里的area是一个拥有*rect类型接收器（receiver）的方法。

可以为值类型或者指针类型的接收者定义方法。代码中展示了一个值类型接收者的例子。

调用方法时，Go会自动处理值和指针之间的转换。想要避免在调用方法时产生一个拷贝，或者想让方法可以修改接受结构体的值，可以使用指针来调用方法。

```Go
package main

import "fmt"

type rect struct {
    width, height int
}

func (r *rect) area() int {
    return r.width * r.height
}

func (r rect) perim() int {
    return 2*r.width + 2*r.height
}

func main() {
    r := rect{width: 10, height: 5}

    fmt.Println("area: ", r.area())
    fmt.Println("perim:", r.perim())

    rp := &r
    fmt.Println("area: ", rp.area())
    fmt.Println("perim:", rp.perim())
}
```

```shell
$ go run methods.go
area:  50
perim: 30
area:  50
perim: 30
```

## 接口

方法签名的集合叫做:`_接口(Interface)_`

代码中创建一个几何体的基本接口。这个例子中，为rect和circle实现该接口。要在Go中实现一个接口，只需要实现接口中的所有方法。这里我们为rect实现了geometry接口。

如果一个变量实现了某个接口，我们就可以调用指定接口中的方法。创建一个通用的measure函数，我们可以通过它来使用所有的geometry。

结构体类型circle和rect都实现了geometry接口，因此可以将其实例作为measure的参数。

```Go
package main

import (
    "fmt"
    "math"
)

type geometry interface {
    area() float64
    perim() float64
}

type rect struct {
    width, height float64
}
type circle struct {
    radius float64
}

func (r rect) area() float64 {
    return r.width * r.height
}
func (r rect) perim() float64 {
    return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
    return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
    return 2 * math.Pi * c.radius
}

func measure(g geometry) {
    fmt.Println(g)
    fmt.Println(g.area())
    fmt.Println(g.perim())
}

func main() {
    r := rect{width: 3, height: 4}
    c := circle{radius: 5}

    measure(r)
    measure(c)
}
```

```shell
$ go run interfaces.go
{3 4}
12
14
{5}
78.53981633974483
31.41592653589793
```

## 嵌入-Embedding

Go支持对于结构体（struct）和接口（interfaces）的嵌入（embedding）以表达一种更加无缝的组合（composition）类型

一个container 嵌入 了一个base，一个嵌入看起来像一个没有名字的字段。

当创建含有嵌入的结构体，必须对嵌入进行显式的初始化，在这里使用嵌入的类型当作字段的名字

可以直接在co上访问base中定义的字段，例如：co.num，或者，也可以使用嵌入的类型名拼写出完整的路径。由于container嵌入了base，因此base的方法也成为了container的方法，在这里直接在co上调用了一个从base嵌入的方法。

可以使用带有方法的嵌入结构来赋予接口实现到其他结构上，因为嵌入了base，在这里可以看到container也实现了describer接口。

```Go
package main

import "fmt"

type base struct {
    num int
}

func (b base) describe() string {
    return fmt.Sprintf("base with num=%v", b.num)
}

type container struct {
    base
    str string
}

func main() {

    co := container{
        base: base{
            num: 1,
        },
        str: "some name",
    }

    fmt.Printf("co={num: %v, str: %v}\n", co.num, co.str)

    fmt.Println("also num:", co.base.num)

    fmt.Println("describe:", co.describe())

    type describer interface {
        describe() string
    }

    var d describer = co
    fmt.Println("describer:", d.describe())
}
```

```shell
$ go run embedding.go
co={num: 1, str: some name}
also num: 1
describe: base with num=1
describer: base with num=1
```

