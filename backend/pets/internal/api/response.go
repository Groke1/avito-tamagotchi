package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status_code int, code string, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(ErrorResponse{Code: code, Message: message})
}
