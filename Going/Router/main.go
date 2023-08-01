package main

import (
	"fmt"
	"going"
	"net/http"
)

func main() {
	r := going.New() // 创建一个路由实例

	r.GET("/", func(c *going.Context) {
		//expect /hello?name=tenet
		c.String(http.StatusOK, "hello %s,you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *going.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *going.Context) {
		c.JSON(http.StatusOK, going.H{"filepath": c.Param("filepath")})
	})

	fmt.Println("Server Running")
	r.Run(":9999")
}
