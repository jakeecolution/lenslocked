package main

import (
	"html/template"
	"os"
)

type UserBase struct {
	Name string
	Age  int
	Bio  string
}

type User struct {
	UserBase
	Address  string
	Animals  []string
	Siblings map[string]UserBase
}

func main() {
	tpl, err := template.ParseFiles("hello.gohtml")

	if err != nil {
		panic(err)
	}
	user := User{
		UserBase: UserBase{
			Name: "Jake",
			Age:  25,
			Bio:  `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		},
		Address:  "123 Main St.",
		Animals:  []string{"Prim", "Cleo", "Bella"},
		Siblings: make(map[string]UserBase),
	}
	user.Siblings["Brother"] = UserBase{
		Name: "John",
		Age:  27,
		Bio:  "My brother",
	}
	user.Siblings["Sister"] = UserBase{
		Name: "Jane",
		Age:  23,
		Bio:  "My sister",
	}
	err = tpl.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
