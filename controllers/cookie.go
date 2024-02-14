package controllers

import (
	"net/http"
	"time"
)

const (
	CookiePrefix    = "lenslocked"
	CookieSessionID = CookiePrefix + "_session_id"
	CookiePath      = "/"
)

func CreateCookieKVP(key, value, path string) *http.Cookie {
	cookie := http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	return &cookie
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := CreateCookieKVP(name, value, CookiePath)
	http.SetCookie(w, cookie)
}

func clearCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   CookiePath,
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

func CreateCookie(key, value, path string, expiration time.Time, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func ReadCookie(key string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
