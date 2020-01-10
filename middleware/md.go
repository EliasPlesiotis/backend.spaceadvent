package middleware

import (
	"html/template"
	"net/http"
)

// Render the pages
func Render(w http.ResponseWriter, name string, data interface{}) {
	t := template.Must(template.ParseGlob("./pages/*.html"))
	if err := t.ExecuteTemplate(w, name, data); err != nil {
		panic(err)
	}
}
