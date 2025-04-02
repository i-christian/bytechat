package server

import "net/http"

func writeError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
