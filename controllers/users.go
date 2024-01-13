package controllers

import (
	"fmt"
	"net/http"

	"github.com/jakeecolution/lenslocked/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
		Name  string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email, password, pwconfirm := r.FormValue("email"), r.FormValue("password"), r.FormValue("password_confirm")
	user := models.NewUser{
		Email:    email,
		Password: password,
	}
	if user.Password != pwconfirm {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}
	myu, err := u.UserService.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<p>User Created: %s </p>", myu.String())
}
