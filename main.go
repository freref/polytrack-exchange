package main

import (
	"context"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
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
		tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "User does not exist"})
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

		if confirm != password {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Passwords don't match"})
			return
		}

		_, err = mail.ParseAddress(email)
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Not a valid email address"})
			return
		}

		var existingUser string
		sql := `SELECT username FROM users WHERE username = $1`
		err = dbpool.QueryRow(context.Background(), sql, username).Scan(&existingUser)
		if err == nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Username already exists"})
			return
		}

		var existingEmail string
		sql = `SELECT username FROM users WHERE email = $1`
		err = dbpool.QueryRow(context.Background(), sql, email).Scan(&existingEmail)
		if err == nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Email already exists"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "We had issues processing your password, please try again"})
		}

		sql = `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

		_, err = dbpool.Exec(context.Background(), sql, username, email, string(hash))
		if err != nil {
			fmt.Println("error insert")
			tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Server failed, please try again"})
		}
		tmpl.ExecuteTemplate(w, "register-error-block", Error{ErrorMessage: "Success"})
	}

	hNav := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/components/nav.html")
	}

	fmt.Printf("main \n")

	http.HandleFunc("/", h1)
	http.HandleFunc("/nav", hNav)
	http.HandleFunc("/login", hLogin)
	http.HandleFunc("/login/submit", hLoginSubmit)
	http.HandleFunc("/register", hRegister)
	http.HandleFunc("/register/submit", hRegisterSubmit)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
