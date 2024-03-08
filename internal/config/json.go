package config

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	// Set content type to JSON and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	// create error payload
	payload := map[string]string{"error": msg}
	// call respondWithJSON
	RespondWithJSON(w, code, payload)
}
