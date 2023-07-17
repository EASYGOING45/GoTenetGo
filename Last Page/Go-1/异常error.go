package main

import (
	"errors"
	"fmt"
)

type user struct {
	name     string
	password string
}

func findUser(user []user, name string) (v *user, err error) {
	for _, u := range user {
		if u.name == name {
			return &u, nil
		}
	} //for_,u 中的_是下标，u是值 nil是空值
	return nil, errors.New("user not found")
}

func main() {
	u, err := findUser([]user{{"wang", "1024"}}, "wang")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u.name) // wang

	if u, err := findUser([]user{{"wang", "1024"}}, "li"); err != nil {
		fmt.Println(err) // not found
		return
	} else {
		fmt.Println(u.name)
	}
}
