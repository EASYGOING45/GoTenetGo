package main

import "fmt"

type user struct {
	name string
	pwd  string
}

//函数定义为func(u user) name() 是将函数定义为结构体的方法
func (u user) checkPassword(password string) bool {
	return u.pwd == password
}

func (u *user) resetPassword(password string) {
	u.pwd = password
}

func main() {
	a := user{name: "John", pwd: "123"}
	fmt.Println(a)
	a.resetPassword("456")

	fmt.Println(a.checkPassword("456"))
}
