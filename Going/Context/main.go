package main

import (
	"going"
	"net/http"
)

func main() {
	r := going.New()
	r.GET("/", func(c *going.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Going FrameWork</h1>")
	})

	r.GET("/hello", func(c *going.Context) {
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
