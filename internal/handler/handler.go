package handler

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/alvinmatias69/wedding_invitation/internal/constant"
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
		handleResponse(w,
			map[string]string{"wWW-authenticate": "Basic"},
			http.StatusUnauthorized,
			[]byte("try ping"))
		return
	}

	if err := h.controller.GetHiddenImage(r.Context(), w); err != nil {
		log.Printf("error while getting hidden image: %v", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
	}
}

func (h *Handler) GetSteamToken(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) == 0 {
		handleResponse(w,
			map[string]string{"wWW-authenticate": "Bearer"},
			http.StatusUnauthorized,
			[]byte("Please take a look more thoroughly"))
		return
	}

	bearerTokens := strings.Split(bearerToken, " ")
	if len(bearerTokens) < 2 {
		handleResponse(w,
			map[string]string{"wWW-authenticate": "Bearer"},
			http.StatusUnauthorized,
			[]byte("Please take a look more thoroughly"))
		return
	}

	res, err := h.controller.GetSteamToken(r.Context(), bearerTokens[1])
	if errors.Is(err, constant.ErrTokenExp) {
		log.Println("token expired")
		handleResponse(w, nil, http.StatusUnauthorized, []byte("your token is already expired please generate a new one"))
		return
	}

	if err != nil {
		log.Printf("error while getting steam token: %v\n", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
		return
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("error while parsing response: %v\n", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
		return
	}

	handleResponse(w, map[string]string{"Content-Type": "application/json"}, http.StatusOK, jsonRes)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		log.Printf("error parsing page query params: %v\n", err)
		page = 0
	}

	messages, err := h.controller.GetMessages(r.Context(), page)
	if err != nil {
		log.Printf("error while getting messages: %v\n", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
		return
	}

	res := entities.GetMessageResponse{
		Messages: messages,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("error while parsing response: %v\n", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
		return
	}

	handleResponse(w, map[string]string{"Content-Type": "application/json"}, http.StatusOK, jsonRes)
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var message entities.Message
	err := dec.Decode(&message)
	if err != nil {
		log.Printf("Error while parsing message: %v\n", err)
		handleResponse(w, nil, http.StatusBadRequest, []byte("error while handling input"))
		return
	}

	err = h.controller.PostMessage(r.Context(), message)
	if err != nil {
		log.Printf("error while parsing response: %v\n", err)
		handleResponse(w, nil, http.StatusInternalServerError, []byte("please contact site admin"))
		return
	}

	handleResponse(w, nil, http.StatusOK, nil)
}
