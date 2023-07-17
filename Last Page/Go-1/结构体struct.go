package main

import "fmt"

type user struct {
	name     string
	password string
}

func main() {
	a := user{name: "LYX", password: "0912"}
	fmt.Println(a)

	var d user
	d.name = "TENET"
	d.password = "tenet"

	fmt.Println(d)

	fmt.Println(checkpasswords(a, "0912"))
	fmt.Println(checkpasswords2(&a, "0912"))
}

func checkpasswords(u user, password string) bool {
	return u.password == password
}

func checkpasswords2(u *user, password string) bool {
	return u.password == password
}
