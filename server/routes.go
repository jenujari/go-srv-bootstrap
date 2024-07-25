package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetRoutes(r *chi.Mux) {
	r.Get("/", Index)
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World 2"))

	if err != nil {
		return
	}
}
