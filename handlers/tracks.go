package handlers

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Track struct {
	Id          int
	Title       string
	Description string
	Code        string
	Vote        int
}

func GetTracks(dbpool *pgxpool.Pool, tracks *[]Track) {
	sql := `SELECT * FROM tracks`
	rows, err := dbpool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track

		err := rows.Scan(&track.Id, &track.Title, &track.Description, &track.Code, &track.Vote)
		if err != nil {
			fmt.Println(err)
			// handle error
		}
		*tracks = append(*tracks, track)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		// handle error
	}
}

func Tracks(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
		var tracks []Track
		GetTracks(dbpool, &tracks)
		tmpl.Execute(w, tracks)
	}
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

		sql := `INSERT INTO tracks (title, code, track_description, vote) VALUES ($1, $2, $3, $4)`

		_, err = dbpool.Exec(context.Background(), sql, title, description, code, 1)
		if err != nil {
			// handle error
		}

		var tracks []Track
		GetTracks(dbpool, &tracks)

		tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
		tmpl.Execute(w, tracks)
	}
}
