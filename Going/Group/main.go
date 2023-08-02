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
