package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jesperkha/notifier"
	"github.com/jesperkha/pipoker/config"
)

type Server struct {
	mux    *chi.Mux
	config *config.Config
}

func New(config *config.Config) *Server {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to PiPoker!"))
	})

	return &Server{
		mux:    mux,
		config: config,
	}
}

func (s *Server) ListenAndServe(notif *notifier.Notifier) {
	done, finish := notif.Register()

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.mux,
	}

	go func() {
		<-done
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		finish()
	}()

	log.Println("listening on port " + s.config.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
	}
}
