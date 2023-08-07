package main

import (
	"fmt"
	"going"
	"net/http"
)

func main() {
	r := going.Default()
	r.GET("/", func(c *going.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *going.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	fmt.Println("Server running........")
	r.Run(":9999")
}
