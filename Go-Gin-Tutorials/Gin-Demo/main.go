// main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//参数路由
	// 无参数
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Who are you?")
	})

	//解析路径参数
	//匹配 /user/tenet
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	//获取Query参数
	//匹配users?name=xxx&role=xxx role为可选字段
	r.GET("/users", func(c *gin.Context) {
		name := c.Query("name")
		role := c.DefaultQuery("role", "teacher")
		c.String(http.StatusOK, "%s is a %s", name, role)
	})

	r.Run(":8181") // listen and serve on 0.0.0.0:8080
}
