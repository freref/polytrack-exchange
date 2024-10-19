package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"polytrack-explorer/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Error struct {
	ErrorMessage string
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		os.Exit(1)
	}
	defer dbpool.Close()

	// components
	http.HandleFunc("/nav", handlers.Nav)
	// pages
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/tracks", handlers.Tracks)
	http.HandleFunc("/tracks/add", handlers.AddTrack)
	// auth
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/login/submit", handlers.LoginSubtmit(dbpool))
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/register/submit", handlers.RegisterSubmit(dbpool))
	// static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
