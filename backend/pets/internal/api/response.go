package api

import (
	"encoding/json"
	"net/http"
)

type LeaderboardItemResponse struct {
	Rank     int    `json:"rank"`
	UserName string `json:"user_name"`
	PetName  string `json:"pet_name"`
	Level    int    `json:"level"`
}

type LeaderboardResponse struct {
	Items       []LeaderboardItemResponse `json:"items"`
	CurrentUser LeaderboardItemResponse   `json:"current_user"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, err APIError) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Code: err.Code, Message: err.Message})
}

<<<<<<< HEAD
func writeJsonResponse(w http.ResponseWriter, status_code int, resp interface{}) {
=======
func writeJsonResponse(w http.ResponseWriter, status_code int, resp PetResponse) {
>>>>>>> 9f0afb9c68d0604e731ec3d40cd30366c4e2a04f
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(resp)
}
