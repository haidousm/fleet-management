package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/html/pages/map.html")
}
