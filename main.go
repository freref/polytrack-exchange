package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		data := Error{
			ErrorMessage: "User does not exist",
		}

		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.ExecuteTemplate(w, "login-error-block", data)
	}

	hRegister := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, nil)
	}

	hRegisterSubmit := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("hRegisterSubmit \n")

		err := r.ParseForm()
		if err != nil {
			fmt.Printf("ERROR parse")
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("ERROR password")
		}

		fmt.Println(username)
		fmt.Println(email)
		fmt.Println(password)
		fmt.Println(string(hash))

		sql := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

		_, err = dbpool.Exec(context.Background(), sql, username, email, string(hash))
		if err != nil {
			fmt.Printf("ERROR inserting user into DB: %v\n", err)
		}

		var res string
		_ = dbpool.QueryRow(context.Background(), "SELECT * FROM users").Scan(&res)

		fmt.Println(res)
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
