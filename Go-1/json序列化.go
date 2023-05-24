package main

import (
	"encoding/json"
	"fmt"
)

type userInfo struct {
	Name  string
	Age   int `json:"age"` // 通过指定tag实现json序列化该字段时的key
	Hobby []string
}

func main() {
	a := userInfo{Name: "wang", Age: 18, Hobby: []string{"篮球", "足球"}}
	buf, err := json.Marshal(a) //Marshal函数将结构体序列化成json字符串

	if err != nil {
		panic(err) //panic函数用于抛出异常
	}

	fmt.Println(buf) //
	fmt.Println(string(buf))

	buf, err = json.MarshalIndent(a, "", "\t") //MarshalIndent函数将结构体序列化成带缩进格式的json字符串
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))

	var b userInfo
	err = json.Unmarshal(buf, &b) //Unmarshal函数将json字符串反序列化成结构体

	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", b) //%#v格式化输出结构体
}
