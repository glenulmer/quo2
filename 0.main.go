package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"

	. "pm/lib/output"
)

func main() {
	Bootstrap()
	defer App.DB.Close()

	r := chi.NewRouter()
	r.Get(`/`, Page0Home)

	r.Handle(`/static/*`, http.StripPrefix(`/static/`, http.FileServer(http.Dir(`./static`))))

	Log(`quo2 is running on :`, App.port)
	log.Fatal(http.ListenAndServe(Str(`:`, App.port), r))
}
