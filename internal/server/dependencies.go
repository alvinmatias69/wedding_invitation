package server

import "net/http"

type handler interface {
	GetHiddenImage(http.ResponseWriter, *http.Request)
	GetSteamToken(http.ResponseWriter, *http.Request)
	GetMessages(http.ResponseWriter, *http.Request)
	PostMessage(http.ResponseWriter, *http.Request)
}
