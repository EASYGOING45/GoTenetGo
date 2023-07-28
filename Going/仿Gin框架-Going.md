# 仿Gin框架-Going

## Day0-框架设计

大部分时候，我们需要实现一个 Web 应用，第一反应是应该使用哪个框架。不同的框架设计理念和提供的功能有很大的差别。比如 Python 语言的 `django`和`flask`，前者大而全，后者小而美。Go语言/golang 也是如此，新框架层出不穷，比如`Beego`，`Gin`，`Iris`等。那为什么不直接使用标准库，而必须使用框架呢？在设计一个框架之前，我们需要回答框架核心为我们解决了什么问题。只有理解了这一点，才能想明白我们需要在框架中实现什么功能。

我们先看看标准库`net/http`如何处理一个请求。

```Go
func main(){
    http.HandleFunc("/",handler)
    http>HandleFunc("/count",counter)
    log.Fatal(http.ListenAndServe("localhost:8080",nil))
}

func handler(w http.ResponseWriter,r *http.Request){
    fmt.Fprintf(w,"URL.Path = %q\n",r.URL.Path)
}
```

`net/http`提供了基础的Web功能，即监听端口，映射静态路由，解析HTTP报文，一些Web开发中简单的需求并不支持，需要手工实现

- 动态路由：例如`hello/:name`，`hello/*`这类的规则。
- 鉴权：没有分组/统一鉴权的能力，需要在每个路由映射的handler中实现。
- 模板：没有统一简化的HTML机制。
- …

当离开框架，使用基础库时，需要频繁手工处理的地方，就是框架的价值所在。但并不是每一个频繁处理的地方都适合在框架中完成。Python有一个很著名的Web框架，名叫[`bottle`](https://github.com/bottlepy/bottle)，整个框架由`bottle.py`一个文件构成，共4400行，可以说是一个微框架。那么理解这个微框架提供的特性，可以帮助我们理解框架的核心能力。

- 路由(Routing)：将请求映射到函数，支持动态路由。例如`'/hello/:name`。
- 模板(Templates)：使用内置模板引擎提供模板渲染机制。
- 工具集(Utilites)：提供对 cookies，headers 等处理机制。
- 插件(Plugin)：Bottle本身功能有限，但提供了插件机制。可以选择安装到全局，也可以只针对某几个路由生效。
- …

## HTTP基础

### 标准库启动Web服务

Go语言内置了 `net/http`库，封装了HTTP网络编程的基础的接口，我们实现的`Gee` Web 框架便是基于`net/http`的。我们接下来通过一个例子，简单介绍下这个库的使用。

```Go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server Running...")
	http.HandleFunc("/", indexHandler)           //设置访问路由
	http.HandleFunc("/hello", helloHandler)      //设置访问路由
	log.Fatal(http.ListenAndServe(":4826", nil)) //设置监听端口
}

// 处理器函数 handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// 处理器函数 handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

```

我们设置了2个路由，`/`和`/hello`，分别绑定 *indexHandler* 和 *helloHandler* ， 根据不同的HTTP请求会调用不同的处理函数。访问`/`，响应是`URL.Path = /`，而`/hello`的响应则是请求头(header)中的键值对信息。

用 curl 这个工具测试一下，将会得到如下的结果。

```
$ curl http://localhost:9999/
URL.Path = "/"
$ curl http://localhost:9999/hello
Header["Accept"] = ["*/*"]
Header["User-Agent"] = ["curl/7.54.0"]
```

*main* 函数的最后一行，是用来启动 Web 服务的，第一个参数是地址，`:9999`表示在 *9999* 端口监听。而第二个参数则代表处理所有的HTTP请求的实例，`nil` 代表使用标准库中的实例处理。第二个参数，则是我们基于`net/http`标准库实现Web框架的入口。

![image-20230728153844790](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230728153844790.png)

### 实现http.Handler接口

```go
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```

第二个参数的类型是什么呢？通过查看`net/http`的源码可以发现，`Handler`是一个接口，需要实现方法 *ServeHTTP* ，也就是说，只要传入任何实现了 *ServerHTTP* 接口的实例，所有的HTTP请求，就都交给了该实例处理了。

```Go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	fmt.Println("Base 2 Server Running On Port 4826...")
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":4826", engine))
}

```

- 定义了一个空的结构体`Engine`，实现了方法`ServeHTTP`。这个方法有2个参数，第二个参数是 *Request* ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；第一个参数是 *ResponseWriter* ，利用 *ResponseWriter* 可以构造针对该请求的响应。
- 在 *main* 函数中，我们给 *ListenAndServe* 方法的第二个参数传入了刚才创建的`engine`实例。至此，我们走出了实现Web框架的第一步，即，将所有的HTTP请求转向了我们自己的处理逻辑。还记得吗，在实现`Engine`之前，我们调用 *http.HandleFunc* 实现了路由和Handler的映射，也就是只能针对具体的路由写处理逻辑。比如`/hello`。但是在实现`Engine`之后，我们拦截了所有的HTTP请求，拥有了统一的控制入口。在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志、异常处理等。

![image-20230728154717477](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230728154717477.png)

### Going框架雏形

重新组织上述代码，搭建`Going`框架的出行

代码目录结构如下：

![image-20230728161638654](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230728161638654.png)

#### base3/go.mod

```go
module example

go 1.19

require going v0.0.0

replace going => ./going

```

#### base3/main.go

```Go
package main

import (
	"fmt"
	"net/http"

	"going"
)

func main() {
	r := going.New() //创建一个going实例
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":4826")
}

```

看到这里，如果你使用过`gin`框架的话，肯定会觉得无比的亲切。`gee`框架的设计以及API均参考了`gin`。使用`New()`创建 gee 的实例，使用 `GET()`方法添加路由，最后使用`Run()`启动Web服务。这里的路由，只是静态路由，不支持`/hello/:name`这样的动态路由，动态路由我们将在下一次实现。

#### base3/going/go.mod

```Go
module going

go 1.19

```



#### base3/going/going.go

```Go
package going

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by going
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of going.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern //将请求方法和路由规则组合成一个唯一的key
	engine.router[key] = handler  //将路由和处理函数绑定
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	fmt.Println("Server is running at " + addr)
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

```

那么`going.go`就是重头戏了。我们重点介绍一下这部分的实现。

- 首先定义了类型`HandlerFunc`，这是提供给框架用户的，用来定义路由映射的处理方法。我们在`Engine`中，添加了一张路由映射表`router`，key 由请求方法和静态路由地址构成，例如`GET-/`、`GET-/hello`、`POST-/hello`，这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法。
- 当用户调用`(*Engine).GET()`方法时，会将路由和处理方法注册到映射表 *router* 中，`(*Engine).Run()`方法，是 *ListenAndServe* 的包装。
- `Engine`实现的 *ServeHTTP* 方法的作用就是，解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 *404 NOT FOUND* 。