package handlers

import (
	"net/http"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Leaderboards(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/leaderboards.html"))
		tmpl.Execute(w, nil)
	}
}
