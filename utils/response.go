package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Success bool `json:"success"`
	Error string  `json:"error"`
	Details string `json:"details"`
	Data interface{} `json:"data,omitempty"`
}

func RespondWithError(w http.ResponseWriter, code int, message string, details string ) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Error: message,
		Details: details,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	response := ApiResponse{
		Success: true,
		Data:    payload,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}