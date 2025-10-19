package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string  `json:"error"`
	Details string `json:"details"`
}

func RespondWithError(w http.ResponseWriter, code int, message string, details string ) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Details: details,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}