package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Track struct {
	Id          int
	Title       string
	Description string
	Code        string
}

func GetTracks(dbpool *pgxpool.Pool) []Track {
	var tracks []Track
	sql := `SELECT * FROM tracks ORDER BY Id DESC`
	rows, err := dbpool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err)
		return tracks
	}
	defer rows.Close()

	for rows.Next() {
		var track Track

		err := rows.Scan(&track.Id, &track.Title, &track.Description, &track.Code)
		if err != nil {
			fmt.Println(err)
			// handle error
			return tracks
		}
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		// handle error
	}

	return tracks
}

func Tracks(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
		tracks := GetTracks(dbpool)
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
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Server error"})
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		code := r.FormValue("code")

		if strings.Contains(code, " ") {
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Invalid track code"})
			return
		}

		if len(title) > 255 {
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Title too long"})
			return
		}

		if len(description) > 1000 {
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Description too long"})
			return
		}

		var existingCode string
		sql := `SELECT code FROM tracks WHERE code = $1`
		err = dbpool.QueryRow(context.Background(), sql, code).Scan(&existingCode)
		if err == nil {
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Track with this code has already been uploaded"})
			return
		}

		sql = `INSERT INTO tracks (title, track_description, code) VALUES ($1, $2, $3)`
		_, err = dbpool.Exec(context.Background(), sql, title, description, code)
		if err != nil {
			tmpl := template.Must(template.ParseFiles("./templates/components/track-form.html"))
			tmpl.Execute(w, Error{ErrorMessage: "Database error, try again"})
		}

		tracks := GetTracks(dbpool)

		tmpl := template.Must(template.ParseFiles("./templates/tracks.html"))
		tmpl.Execute(w, tracks)
	}
}
