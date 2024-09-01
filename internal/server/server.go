package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
)

type Server struct {
	handler handler
}

func New(handler handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Start(cfg entities.Config) {
	fs := http.FileServer(http.Dir(cfg.StaticWebDir))
	http.Handle("GET /", fs)

	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc(fmt.Sprintf("GET %s", cfg.HiddenImagePath), s.handler.GetHiddenImage)
	http.HandleFunc(fmt.Sprintf("GET %s", cfg.FinalImagePath), s.handler.GetFinalImage)

	fmt.Printf("starting server in port: %v\n", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
