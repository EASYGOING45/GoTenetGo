package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"going"
)

func onlyForV2() going.HandlerFunc {
	return func(c *going.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := going.New()
	r.Use(going.Logger()) // global midlleware
	r.GET("/", func(c *going.Context) {
		c.HTML(http.StatusOK, "<h1>Hello going</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *going.Context) {
			// expect /hello/goingktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	fmt.Println("Server running at :9999")
	r.Run(":9999")
}
