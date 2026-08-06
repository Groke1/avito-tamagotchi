package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, err APIError) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Code: err.Code, Message: err.Message})
}

func writeJsonResponse(w http.ResponseWriter, status_code int, resp PetResponse) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(resp)
}
