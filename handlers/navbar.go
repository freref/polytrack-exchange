package handlers

import "net/http"

func Navbar(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/components/navbar.html")
}
