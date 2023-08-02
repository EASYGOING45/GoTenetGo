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

	Params map[string]string // 路由参数

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

// Param方法可以获取路由参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
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
