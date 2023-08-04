# Going-PartⅥ-模板HTML Template

- 实现静态资源服务(Static Resource)
- 支持HTML模板渲染

## 服务端渲染

目前越来越流行前后端分离的开发模式，即Web后端提供RESTful接口，返回结构化的数据(通常为JSON或者XML)。前端使用AJAX技术请求到所需的数据，利用JavaScript进行渲染。

VUE/React等前端框架持续火热，这样的开发模式前后端解耦，优势突出，后端只需专心解决资源利用、并发、数据库等问题，只需考虑数据如何生成；前端专注于界面的设计实现，只需考虑拿到数据后如何渲染即可。

前后端分离还有一大优势，因为后端只关注于数据，接口返回值是结构化的，与前端解耦，同一套后端服务能够同时支撑小程序、移动App、PC端H5等，以及对外提供的接口，随着前端工程化的不断发展，Webpack，gulp等工具层出不穷，前端技术越发自成体系。

但前后分离的一大问题在于，页面是在客户端渲染的，比如浏览器，这对于爬虫并不友好。Google爬虫已经能够爬取渲染后的网页，但是短期内爬取服务端直接渲染的HTML页面仍是主流。

今日主要实现Web框架如何支持服务端渲染的场景。

## 静态文件 Serve Static Files

网页三剑客，JavaScript、CSS和HTML，要做到服务端渲染，第一步便是要支持JS、CSS等静态文件，记得之前设计动态路由的时候，支持通配符 `*`匹配多级子路径。比如路由 `/assets/*filepath`，可以匹配 `/assets/`开头的所有的地址。例如 `/assets/js/ctenet.js`，匹配后，参数 `filepath`就赋值为 `js/ctenet.js`

那如果将所有的静态文件都放在 `/usr/web`目录下，那么 `filepath`的值即是该目录下文件的相对地址，映射到真实的文件后，将文件返回，静态服务器就是先了。

找到文件后，如何返回这一步， `net/http`库已经实现了，因此，going框架要做的，仅仅是解析请求的地址，映射到服务器上文件的真实地址，交给 `http.FileServer`处理就好了。

```Go
// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}
```

给`RouterGroup`添加了2个方法，`Static`这个方法是暴露给用户的。用户可以将磁盘上的某个文件夹`root`映射到路由`relativePath`。例如:

```Go
r := gee.New()
r.Static("/assets", "/usr/geektutu/blog/static")
// 或相对路径 r.Static("/assets", "./static")
r.Run(":9999")
```

用户访问`localhost:9999/assets/js/geektutu.js`，最终返回`/usr/geektutu/blog/static/js/geektutu.js`。

## HTML模板渲染

Go语言内置了`text/template`和`html/template`2个模板标准库，其中[html/template](https://golang.org/pkg/html/template/)为 HTML 提供了较为完整的支持。包括普通变量渲染、列表渲染、对象渲染等。gee 框架的模板渲染直接使用了`html/template`提供的能力。

```Go
Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
```

首先为 Engine 示例添加了 `*template.Template` 和 `template.FuncMap`对象，前者将所有的模板加载进内存，后者是所有的自定义模板渲染函数。

另外，给用户分别提供了设置自定义渲染函数`funcMap`和加载模板的方法。

接下来，对原来的 `(*Context).HTML()`方法做了些小修改，使之支持根据模板文件名选择模板进行渲染。

```Go
type Context struct {
    // ...
	// engine pointer
	engine *Engine
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
```

我们在 `Context` 中添加了成员变量 `engine *Engine`，这样就能够通过 Context 访问 Engine 中的 HTML 模板。实例化 Context 时，还需要给 `c.engine` 赋值。

```Go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// ...
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
```

```
---gee/
---static/
   |---css/
        |---geektutu.css
   |---file1.txt
---templates/
   |---arr.tmpl
   |---css.tmpl
   |---custom_func.tmpl
---main.go
```

```html
<!-- day6-template/templates/css.tmpl -->
<html>
    <link rel="stylesheet" href="/assets/css/geektutu.css">
    <p>geektutu.css is loaded</p>
</html>
```

```Go
type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
```

![image-20230804141529002](https://happygoing.oss-cn-beijing.aliyuncs.com/img/image-20230804141529002.png)

