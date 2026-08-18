package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encoding response: %v", err)
	}
}

func OKResponse(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusOK, data)
}

func CreatedResponse(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusCreated, data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}
