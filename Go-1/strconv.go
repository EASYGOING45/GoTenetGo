package main

import (
	"fmt"
	"strconv"
)

// strconv包实现了基本数据类型与其字符串表示的相互转换
func main() {
	f, _ := strconv.ParseFloat("1.234", 64)
	fmt.Println(f) // 1.234

	n, _ := strconv.ParseInt("111", 10, 64)
	fmt.Println(n) // 111

	n, _ = strconv.ParseInt("0x1000", 0, 64)
	fmt.Println(n) // 4096

	n2, _ := strconv.Atoi("123")
	fmt.Println(n2) // 123

	n2, err := strconv.Atoi("AAA") //Atoi
	fmt.Println(n2, err)           // 0 strconv.Atoi: parsing "AAA": invalid syntax
}
