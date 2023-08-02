# Going-PartⅣ-分组控制Group

- 实现路由的分组控制(Route Group Control)

## 分组的意义

分组控制(Group Control)是Web框架应提供的基础功能之一。所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制。但是在真实的业务场景中，往往某一组路由需要相似的处理。例如：

- 以 `/post`开头的路由匿名可访问
- 以 `/admin`开头的路由需要鉴权
- 以 `/api`开头的路由是RESTful接口，可以对接第三方平台，需要三方平台鉴权

大部分情况下的路由分组，是以相同的前缀来区分的。因此这里实现的分组控制也是以前缀来区分，并且支持分组的嵌套。例如 `/post`是一个分组，`/post/a`和 `/post/b`可以是该分组下的子分组。作用在 `/post`分组上的中间件(middleware)，也都会作用在子分组，子分组还可以应用自己特有的中间件

中间件可以给框架提供无限的扩展能力，应用在分组上，可以使得分组控制的收益更为明显，而不是共享相同的路由前缀这么简单，例如 `/admin`的分组，可以应用鉴权中间件；`/`分组应用日志中间件，`/`是默认的最顶层的分组，也就意味着给所有的路由，即整个框架增加了记录日志的能力。

提供扩展能力和支持中间件的内容，将在下一部分进行。

## 分组嵌套

一个Group对象需要具备哪些属性呢？首先是前缀(prefix)，例如 `/`，或者 `/api`，要支持分组嵌套，那么需要知道当前分组的父亲是谁(parent)，当然，按照一开始的分析，中间件是应用在分组上的，那还需要存储应用在该分组上的中间件(middlewares)，还记得，之前调用函数 `(*Engine).addRoute()`来映射所有的路由规则和Handler。如果Group对象需要直接映射路由规则的话，比如想在使用框架时，可以这么调用:

```Go
r := gee.New()
v1 := r.Group("/v1")
v1.GET("/", func(c *gee.Context) {
	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
})
```

那么Group对象，还需要有访问 `Router`的能力，为了方便，我们可以在Group中，保存一个指针，指向Engine，整个框架的所有资源都是由 `Engine`统一协调的，那么就可以通过 `Engine`间接地访问各种接口了。

所以，最后的Group的定义是这样的

### going.go

```Go
RouterGroup struct{
    prefix 		string
    middlewares	[]HandlerFunc	//支持中间件
    parent		*RouterGroup	//Support nesting
    engine		*Engine			//共享引擎初始化
}
```

我们还可以进一步地抽象，将`Engine`作为最顶层的分组，也就是说`Engine`拥有`RouterGroup`所有的能力

```Go
Engine struct{
    *RouterGroup
    router	*router
    groups  []*RouterGroup	//存储所有的Groups
}
```

接下来就可以将所有和路由相关的函数都交给`RouterGroup`来实现了

```Go
// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
```

可以仔细观察下`addRoute`函数，调用了`group.engine.router.addRoute`来实现了路由的映射，由于`Engine`从某种意义上继承了`RouterGroup`的所有属性和方法，因为`(*Engine).engines`是指向自己的，这样实现，我们既可以像原来一样添加路由，也可以通过分组来添加路由。

## Demo

```Go
package main

import (
	"fmt"
	"going"
	"net/http"
)

func main() {
	r := going.New()
	r.GET("/index", func(c *going.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *going.Context) {
			c.HTML(http.StatusOK, "<h1>Hello v1</h1>")
		})

		v1.GET("/hello", func(c *going.Context) {
			c.String(http.StatusOK, "hello %s,you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *going.Context) {
			// expect /hello/goingktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *going.Context) {
			c.JSON(http.StatusOK, going.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	fmt.Println("Server is running at http://localhost:9999")
	r.Run(":9999")
}

```

![image-20230802103821443](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230802103821443.png)