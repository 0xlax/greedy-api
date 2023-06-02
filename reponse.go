package main

import (
	"encoding/json"
	"net/http"
)

// Sends error response to the client.
func sendErrorResponse(w http.ResponseWriter, errorMessage string) {
	// Create ErrorResponse object as JSON with the specified error message.
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
}

// Sends a value response.
func sendValueResponse(w http.ResponseWriter, value string) {
	// CreateValueResponse object as JSON with the specified value.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValueResponse{Value: value})
}

// Sends a simple OK response to the client.
func sendOKResponse(w http.ResponseWriter) {
	// Send an empty response as JSON to indicate a successful response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
