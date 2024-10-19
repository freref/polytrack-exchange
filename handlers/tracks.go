package handlers

import (
	"context"
	"net/http"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Tracks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
	tmpl.Execute(w, nil)
}

func AddTrack(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
	tmpl.Execute(w, nil)
}

func SubmitTrack(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// handle error
			return
		}

		title := r.FormValue("title")
		code := r.FormValue("description")
		description := r.FormValue("code")

		sql := `INSERT INTO tracks (title, code, track_description, upvote, downvote) VALUES ($1, $2, $3, $4, $5)`

		_, err = dbpool.Exec(context.Background(), sql, title, description, code, 1, 0)
		if err != nil {
			// handle error
		}
		tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
		tmpl.Execute(w, nil)
	}
}
