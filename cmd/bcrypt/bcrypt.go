package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		if len(os.Args) < 4 {
			fmt.Println("Usage: bcrypt compare [hash] [password]")
			return
		}
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Println("Usage: bcrypt [hash|compare]")
	}
}

func hash(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(hash))
}

func compare(hash, password string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Passwords do not match")
		return
	}
	fmt.Println("Passwords match")
}
