package handlers

import "net/http"

func Nav(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/components/nav.html")
}
