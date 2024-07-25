package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jenujari/go-srv-bootstrap/config"
	"github.com/jenujari/go-srv-bootstrap/helpers"
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
