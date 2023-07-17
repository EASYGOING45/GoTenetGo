package main

import "fmt"

func add2(n int) {
	n += 2 //该方式无法传递正确
}

func add2ptr(n *int) {
	*n += 2
}

func main() {
	var n int = 10
	add2(n)
	fmt.Println(n) // 12

	add2ptr(&n)
	fmt.Println(n) // 14
}
