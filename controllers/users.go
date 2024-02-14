package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jakeecolution/lenslocked/models"
)

type Users struct {
	Templates struct {
		New         Template
		SignIn      Template
		CurrentUser Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(CookieSessionID, r)
	if err != nil {
		// fmt.Printf("%v, %T\n", r, r)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	} else if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong with the signin", http.StatusInternalServerError)
		return
	}
	params := r.URL.Query()
	if params.Get("Success") == "true" {
		var data struct {
			Email    string
			Password string
			Error    string
			Success  bool
		}
		data.Email = user.Email
		data.Password = ""
		data.Error = ""
		data.Success = true
		u.Templates.CurrentUser.Execute(w, r, data)
		return
	}
	var data struct {
		Email    string
		Password string
		Error    string
		Success  bool
	}
	data.Email = user.Email
	data.Password = ""
	data.Error = params.Get("Errors")
	data.Success = false
	fmt.Println(data)
	u.Templates.CurrentUser.Execute(w, r, data)

	// user, err := u.UserService.ByEmail(cookie.Value)
	// if err != nil {
	// 	return nil
	// }
	// return user
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
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
	session, err := u.SessionService.Create(myu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	setCookie(w, CookieSessionID, session.TokenHash)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email, password := r.FormValue("email"), r.FormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusNotAcceptable)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong with the signin", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSessionID, session.TokenHash)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(CookieSessionID, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		http.Error(w, "Something went wrong with the signout", http.StatusInternalServerError)
		return
	}
	clearCookie(w, CookieSessionID)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ProcessUserUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(CookieSessionID, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	} else if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong with the signin", http.StatusInternalServerError)
		return
	}
	var data struct {
		Email           string
		Password        string
		PasswordConfirm string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	data.PasswordConfirm = r.FormValue("password_confirm")
	if data.Password != data.PasswordConfirm {
		burl, err := url.Parse("")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong with the signin", http.StatusInternalServerError)
			return
		}
		burl.Path = "/users/me"
		q := url.Values{}
		q.Add("Errors", "Passwords do not match")
		burl.RawQuery = q.Encode()
		http.Redirect(w, r, burl.String(), http.StatusFound)
		return
	} else if (data.Password == "") && (data.Password == data.PasswordConfirm) {
		fmt.Println("Did not update password")
	} else {
		user.HashedPassword, err = models.HashPass(data.Password)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong with the signin", http.StatusInternalServerError)
			return
		}
	}
	if data.Email != user.Email {
		user.Email = data.Email
	}
	err = u.UserService.Update(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong with updating the user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users/me?Success=true", http.StatusFound)
	// user, err := u.UserService.ByEmail(cookie.Value)
	// if err != nil {
	// 	return nil
	// }
	// return user
}
