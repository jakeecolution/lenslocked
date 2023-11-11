package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, path string) (Template, error) {
	tpl, err := template.ParseFS(fs, path)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %v", err)
	}
	return Template{htmlTpl: tpl}, nil
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %v", err)
	}
	return Template{htmlTpl: tpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully executed template.")
}
