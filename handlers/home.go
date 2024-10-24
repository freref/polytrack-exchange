package handlers

import (
	"net/http"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, nil)
}

func HomeContent(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	_, ok := headers[http.CanonicalHeaderKey("HX-Request")]

	if ok {
		tmpl := template.Must(template.ParseFiles("./templates/home.html"))
		tmpl.Execute(w, nil)
	} else {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	}
}
