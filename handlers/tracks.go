package handlers

import (
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

func AddTrackSubmit(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
