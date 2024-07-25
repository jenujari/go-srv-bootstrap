package server

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jenujari/go-srv-bootstrap/config"
	"github.com/jenujari/go-srv-bootstrap/helpers"
	"github.com/jenujari/go-srv-bootstrap/tpl"
)

var (
	server *http.Server
	router *chi.Mux
)

func init() {

	server = &http.Server{
		Addr:              ":5456",
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	router = chi.NewRouter()
	router.Use(middleware.Logger)

	FileServer(router, "/static")
	SetRoutes(router)

	server.Handler = router
	config.GetLogger().Println("server initialization complete.")
}

func RunServer(cmder *helpers.Commander) {
	defer cmder.CompleteOneWorker()

	go func(cmdx *helpers.Commander) {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			cmdx.FatalErrorChan <- fmt.Errorf("ListenAndServe(): %v", err)
		}
	}(cmder)

	// helper.HitBrowser("http://localhost:5456", j)

	<-cmder.CTX.Done()
	config.GetLogger().Println("shutting down server...")
	if err := server.Shutdown(cmder.CTX); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	config.GetLogger().Println("server shutdown complete...")
}

func GetServer() *http.Server {
	return server
}


// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string) {
	sub, err := fs.Sub(tpl.GetAssetsFs(), "assets")
	if err != nil {
		panic(err)
	}

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, rx *http.Request) {
		rctx := chi.RouteContext(rx.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*") + "/"
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(sub)))
		fs.ServeHTTP(w, rx)
	})
}
