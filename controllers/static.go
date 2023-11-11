package controllers

import (
	"net/http"

	"github.com/jakeecolution/lenslocked/views"
)

type Static struct {
	Template views.Template
}

func StaticHandler(tpl views.Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}
