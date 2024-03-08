package config

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusInternalServerError, "Internal server error")
}
