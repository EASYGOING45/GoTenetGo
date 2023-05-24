package main

import "fmt"

func main() {

	s := make([]string, 3) // make([]T, length, capacity) 用于创建动态数组
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("get:", s[2])   // c
	fmt.Println("len:", len(s)) // 3
	//Println是一个用于输出的函数，它可以接受任意类型的数据作为参数，多个参数之间用空格分隔，输出结果会自动添加空格。

	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println(s) // [a b c d e f]

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println(c) // [a b c d e f]

	fmt.Println(s[2:5]) // [c d e]		切片区间为左闭右开
	fmt.Println(s[:5])  // [a b c d e]
	fmt.Println(s[2:])  // [c d e f]

	good := []string{"g", "o", "o", "d"}
	fmt.Println(good) // [g o o d]
}
