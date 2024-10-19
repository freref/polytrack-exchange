package handlers

import (
	"net/http"
	"text/template"
)

func Maps(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/maps.html"))
	tmpl.Execute(w, nil)
}
