package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "hello"
	//Contains函数判断字符串是否包含子串
	fmt.Println(strings.Contains(a, "ll")) // true Contains函数判断字符串是否包含子串

	//Count函数统计字符串出现的次数
	fmt.Println(strings.Count(a, "l")) // 2

	//HasPrefix函数判断字符串是否以某个字符串开头
	fmt.Println(strings.HasPrefix(a, "he")) // true

	//HasSuffix函数判断字符串是否以某个字符串结尾
	fmt.Println(strings.HasSuffix(a, "llo")) // true

	//Index函数返回字符串第一次出现的位置，如果没有返回-1
	fmt.Println(strings.Index(a, "ll")) // 2

	//LastIndex函数返回字符串最后一次出现的位置，如果没有返回-1
	fmt.Println(strings.Join([]string{"he", "llo"}, "-")) // he-llo

	//Repeat函数将字符串重复n次
	fmt.Println(strings.Repeat(a, 2)) // hellohello

	//Replace函数将字符串中的某个字符串替换成另一个字符串，n表示替换几个，如果n=-1表示全部替换
	fmt.Println(strings.Replace(a, "e", "E", -1)) // hEllo

	//Split函数将字符串以某个字符串为分割标准分割成若干个子串，返回一个切片
	fmt.Println(strings.Split("a-b-c", "-")) // [a b c]

	//ToLower函数将字符串中的大写字母转换成小写字母
	fmt.Println(strings.ToLower(a)) // hello

	//ToUpper函数将字符串中的小写字母转换成大写字母
	fmt.Println(strings.ToUpper(a)) // HELLO

	//Trim函数将字符串前后的某个字符串去掉
	fmt.Println(len(a)) // 5
	b := "你好"
	fmt.Println(len(b)) // 6
}
