package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	}

	hLogin := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
	}

	hLoginSubmit := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
	}

	hRegister := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, nil)
	}

	hNav := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/components/nav.html")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/nav", hNav)
	http.HandleFunc("/login", hLogin)
	http.HandleFunc("/login/submit", hLoginSubmit)
	http.HandleFunc("/register", hRegister)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
