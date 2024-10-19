package handlers

import (
	"context"
	"net/http"
	"net/mail"
	"os"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	ErrorMessage string
}

func createJWT(username string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/login.html"))
	tmpl.Execute(w, nil)
}

func LoginSubtmit(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		tokenString, err := createJWT(username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
		})

		tmpl.ExecuteTemplate(w, "login-error-block", Error{ErrorMessage: "Success"})
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/register.html"))
	tmpl.Execute(w, nil)
}

func RegisterSubmit(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
