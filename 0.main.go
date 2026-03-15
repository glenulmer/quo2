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
	r.Get(`/`, Page2FiltersGet)
	r.Post(postStateReset, Page2StateResetPost)
	r.Post(postCustomerState, Page2CustomerPost)
	r.Post(postFiltersState, Page2FiltersPost)
	r.Handle(`/static/*`, http.StripPrefix(`/static/`, http.FileServer(http.Dir(`./static`))))

	log.Println(`quo2 is running on :`, App.Port)
	log.Fatal(http.ListenAndServe(Str(`:`, App.Port), r))
}
