package main

import "fmt"

func main() {
	nums := []int{2, 3, 4}
	sum := 0
	for i, num := range nums {
		sum += num
		fmt.Println("index:", i, "num:", num)
		fmt.Println("sum:", sum)
	}
	//range是一个迭代器，每次迭代都会返回两个值，第一个是元素的下标，第二个是元素的值
	m := map[string]string{"a": "A", "b": "B"}
	for k, v := range m {
		fmt.Println(k, v) // b 8; a A
	}
	for k := range m {
		fmt.Println("key", k) // key a; key b
	}
}
