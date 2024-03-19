package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetRoutes(r *chi.Mux) {
	r.Get("/", Index)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World 2"))
}
