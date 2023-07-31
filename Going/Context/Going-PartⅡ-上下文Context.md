# Going-PartⅡ-上下文Context

设计编写Web框架-Going的第二天，今天主要任务成果如下：

- 将`路由（router）`独立出来，方便之后扩展增强
- 设计 `上下文（context）`，封装请求和响应 `Request & Response`，提供对JSON、HTML等返回类型的支持

## 效果演示

main.go

```Go
package main

import (
	"going"
	"net/http"
)

func main() {
	r := going.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Going FrameWork</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		//expect /hello?name=ctenet
		c.String(http.StatusOK, "hello %s,you're at %s\n", c.Query("name"), c.Path)

	})

	r.POST("/login", func(c *going.Context) {
		c.JSON(http.StatusOK, going.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":4826")
}
```

1. `Handler`的参数变成了 `going.Context`，提供了查询Query/PostForm参数的功能
2. `going.Context`封装了 `HTML/String/JSON`函数，能够快速构造HTTP响应

## 设计Context

**必要性**

1. 对Web服务来说，无非是根据请求 `*http.Request`，构造响应 `http.ResponseWriter`。单是这两个对象提供的接口粒度太细，比如我们要构造一个完整的响应，需要考虑消息头（header）和消息体（body），而Header包含了状态码（StatusCode），消息类型（ContentType）等几乎每次请求都需要设置的信息。因此，如果不进行有效的封装，那么框架的用户将需要写大量重复，繁杂的代码，而且容易出错。针对常用场景，能够高效地构造出HTTP响应是一个好的框架必须考虑的点

   用返回JSON数据作比较，感受下封装前后的差距

   封装前

   ```Go
   obj = map[string]interface{}{
       "name": "geektutu",
       "password": "1234",
   }
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusOK)
   encoder := json.NewEncoder(w)
   if err := encoder.Encode(obj); err != nil {
       http.Error(w, err.Error(), 500)
   }
   ```

   封装后

   ```Go
   c.JSON(http.StatusOK, going.H{
       "username": c.PostForm("username"),
       "password": c.PostForm("password"),
   })
   ```

2. 针对使用场景，封装 `*http.Request`和 `http.ResponseWriter`的方法，简化相关接口的调用，只是设计Context的原因之一。对于框架来说，还需要支撑额外的功能。例如，将来解析动态路由 `/hello/:name`，参数 `:name`的值放在哪呢？再比如，框架需要支持中间件，那中间件产生的信息放在哪呢？Context随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由Context承载。因此，设计Context结构，扩展性和复杂性留在了内部，而对外简化了接口。路由的处理函数，以及将要实现的中间件，参数都统一使用Context实例，Context就像一次会话的百宝箱，可以找到任何东西

## 具体实现-Context.go

```Go
package going

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H是一个map[string]interface{}类型的简写，用于生成JSON数据
type H map[string]interface{}

type Context struct {
	//origin objects 是原始的对象
	Writer http.ResponseWriter // ResponseWriter是一个接口，定义了响应的基本操作
	Req    *http.Request       // Request是一个结构体，表示客户端的请求

	//request info 请求信息
	Path   string // 请求路径
	Method string // 请求方法

	//response info 响应信息
	StatusCode int // 响应状态码
}

// newContext是Context的构造函数 用于创建一个Context实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm方法可以获取到表单中对应key的value值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query方法可以获取到url中对应key的value值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status方法可以设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader方法可以设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String方法可以返回字符串类型的响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON方法可以返回JSON类型的响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)

	// 如果在编码的过程中发生错误，应该将错误信息写入到HTTP响应中
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data方法可以返回字节类型的响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML方法可以返回HTML类型的响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

```

- 代码初始位置，给 `map[string]interface{}`起了一个别名 `going.H`,构建JSON数据时，显得更简洁
- `Context`目前只包含了 `http.ResponseWriter`和 `*http.Request`，另外提供了对Method和Path这两个常用属性的直接访问
- 提供了访问`Query`和 `PostForm`参数的方法
- 提供了快速构造 `String/Data/JSON/HTML`响应的方法

## 路由Router

将和路由相关的方法和结构提取出来，构成`router.go`文件，方便之后对router的功能进行增强，如提供对动态路由的支持。router的handle方法作了一些细微调整，即handler的参数，变成了Context。

### router.go

```Go
package going

import (
	"net/http"
)

// 定义HandlerFunc类型，用于路由映射
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + c.Path
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND : %s\n", c.Path)
	}
}

```

## 框架入口

### going.go

```Go
package going

import (
	"log"
	"net/http"
)

// 定义一个结构体 context 用来存储请求的信息
type HandlerFunc func(*Context)

// Engine 是 going 的核心结构，实现 ServeHTTP 接口
type Engine struct {
	router *router
}

// 实现 ServeHTTP 接口 使得 Engine 能够响应 HTTP 请求
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
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
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

```

将Router相关的代码独立后，`going.go`变得简介了许多。最重要的是通过实现了ServeHTTP接口，接管了所有的HTTP请求，相比之前，有了细微的调整，在调用router.handle之前，构造了一个Context对象。这个对象目前还很简单，仅仅是包装了原来的两个参数，后续将变得逐渐清强大。

![image-20230731164326519](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230731164326519.png)