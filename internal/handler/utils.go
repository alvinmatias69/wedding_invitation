package handler

import "net/http"

func handleResponse(w http.ResponseWriter, headers map[string]string, statusCode int, body []byte) {
	for key, val := range headers {
		w.Header().Add(key, val)
	}
	w.WriteHeader(statusCode)
	w.Write(body)
}
