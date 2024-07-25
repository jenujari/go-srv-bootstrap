package server

import (
	"net/http"

	"github.com/jenujari/go-srv-bootstrap/tpl"

	"github.com/go-chi/chi/v5"
)

func SetRoutes(r *chi.Mux) {
	r.Get("/", Index)
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := tpl.GetTemplateExecutor()
	err := t.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		return
	}
}
