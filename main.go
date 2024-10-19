package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	ErrorMessage string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	dbURL := os.Getenv("DATABASE_URL")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		os.Exit(1)
	}
	defer dbpool.Close()

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

		err := r.ParseForm()
		if err != nil {
			tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "Invalid inputs, please try again"})
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		var password_hash string
		sql := `SELECT password_hash FROM users WHERE username = $1`
		err = dbpool.QueryRow(context.Background(), sql, username).Scan(&password_hash)
		if err != nil {
			tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "User doesn't exist"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
		if err != nil {
			tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "Password incorect, try again"})
			return
		}

		tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "Success"})
	}

	hRegister := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, nil)
	}

	hRegisterSubmit := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))

		err := r.ParseForm()
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Invalid inputs, please try again"})
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")

		var existingUser string
		sql := `SELECT username FROM users WHERE username = $1`
		err = dbpool.QueryRow(context.Background(), sql, username).Scan(&existingUser)
		if err == nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Username already exists"})
			return
		}

		var existingEmail string
		sql = `SELECT email FROM users WHERE email = $1`
		err = dbpool.QueryRow(context.Background(), sql, email).Scan(&existingEmail)
		if err == nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Email already exists"})
			return
		}

		_, err = mail.ParseAddress(email)
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Not a valid email address"})
			return
		}

		if confirm != password {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Passwords don't match"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "We had issues processing your password, please try again"})
		}

		sql = `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

		_, err = dbpool.Exec(context.Background(), sql, username, email, string(hash))
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Server failed, please try again"})
		}
		tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Success"})
	}

	hNav := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/components/nav.html")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/nav", hNav)
	http.HandleFunc("/login", hLogin)
	http.HandleFunc("/login/submit", hLoginSubmit)
	http.HandleFunc("/register", hRegister)
	http.HandleFunc("/register/submit", hRegisterSubmit)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
