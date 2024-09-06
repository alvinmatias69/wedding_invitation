package server

import "net/http"

type handler interface {
	GetHiddenImage(http.ResponseWriter, *http.Request)
	GetSteamToken(http.ResponseWriter, *http.Request)
}
