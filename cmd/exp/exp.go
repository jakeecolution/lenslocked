package main

import (
	"fmt"
	"os"

	"github.com/jakeecolution/lenslocked/models"
	msql "github.com/jakeecolution/lenslocked/models/sql"
	"github.com/joho/godotenv"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load(".credentials.env")
	CheckErr(err)
	cfg := msql.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		SSLMode:  false,
	}
	db, err := msql.Open(cfg)
	CheckErr(err)
	defer db.Close()

	err = db.Ping()
	CheckErr(err)
	fmt.Println("Connected!")

	us := models.UserService{
		DB: db,
	}
	superuser, superpassword := os.Getenv("SUPERUSER_EMAIL"), os.Getenv("SUPERUSER_PASSWORD")
	user, err := us.Create(&models.NewUser{Email: superuser, Password: superpassword})
	CheckErr(err)
	fmt.Println(user)
}
