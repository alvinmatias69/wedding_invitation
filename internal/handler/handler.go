package handler

import (
	"crypto/sha512"
	"crypto/subtle"
	"log"
	"net/http"
	"strings"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
)

type Handler struct {
	cfg          entities.Config
	usernameHash [64]byte
	passwordHash [64]byte
	controller   controller
}

func New(cfg entities.Config, controller controller) *Handler {
	return &Handler{
		cfg:          cfg,
		usernameHash: sha512.Sum512([]byte(cfg.Username)),
		passwordHash: sha512.Sum512([]byte(cfg.Password)),
		controller:   controller,
	}
}

func (h *Handler) GetHiddenImage(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("wWW-authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("almost there"))
		return
	}

	var (
		usernameHash  = sha512.Sum512([]byte(username))
		passwordHash  = sha512.Sum512([]byte(password))
		usernameMatch = subtle.ConstantTimeCompare(usernameHash[:], h.usernameHash[:]) == 1
		passwordMatch = subtle.ConstantTimeCompare(passwordHash[:], h.passwordHash[:]) == 1
	)

	if !usernameMatch || !passwordMatch {
		// TODO: move this to reusable function
		w.Header().Add("wWW-authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("almost there"))
		return
	}

	if err := h.controller.GetHiddenImage(r.Context(), w); err != nil {
		log.Printf("error while getting hidden image: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("please contact site administrator"))
	}
}

func (h *Handler) GetFinalImage(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) == 0 {
		// TODO: move this to reusable function
		w.Header().Add("wWW-authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("almost there"))
		return
	}

	bearerTokens := strings.Split(bearerToken, " ")
	if len(bearerTokens) < 2 {
		// TODO: move this to reusable function
		w.Header().Add("wWW-authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("almost there"))
		return
	}

	token := bearerTokens[1]
}
