package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
)

type Server struct {
	handler handler
	cfg     entities.Config
}

func New(cfg entities.Config, handler handler) *Server {
	return &Server{
		handler: handler,
		cfg:     cfg,
	}
}

func (s *Server) Start() {
	fs := http.FileServer(http.Dir(s.cfg.StaticWebDir))
	http.Handle("GET /", fs)

	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc(fmt.Sprintf("GET %s", s.cfg.HiddenImagePath), s.handler.GetHiddenImage)
	http.HandleFunc(fmt.Sprintf("GET %s", s.cfg.SteamTokenPath), s.handler.GetSteamToken)

	fmt.Printf("starting server in port: %v\n", s.cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
