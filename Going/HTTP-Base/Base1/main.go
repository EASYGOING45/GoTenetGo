package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server Running...")
	http.HandleFunc("/", indexHandler)           //设置访问路由
	http.HandleFunc("/hello", helloHandler)      //设置访问路由
	log.Fatal(http.ListenAndServe(":4826", nil)) //设置监听端口
}

// 处理器函数 handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// 处理器函数 handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
