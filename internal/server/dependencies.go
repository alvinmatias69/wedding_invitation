package server

import "net/http"

type handler interface {
	GetHiddenImage(http.ResponseWriter, *http.Request)
	GetFinalImage(http.ResponseWriter, *http.Request)
}
