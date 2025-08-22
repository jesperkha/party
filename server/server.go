package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jesperkha/notifier"
	"github.com/jesperkha/pipoker/config"
	"github.com/jesperkha/pipoker/ws"
)

type Server struct {
	ws     *ws.Server
	mux    *chi.Mux
	config *config.Config
}

func New(config *config.Config) *Server {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	ws := ws.NewServer()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	mux.Get("/client.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/client.js")
	})

	mux.Get("/connect", ws.ConnectHandler())

	return &Server{
		ws,
		mux,
		config,
	}
}

func (s *Server) ListenAndServe(notif *notifier.Notifier) {
	done, finish := notif.Register()

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.mux,
	}

	go s.ws.Run(notif)

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
