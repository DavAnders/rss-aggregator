package config

import (
	"net/http"
	"strings"

	"github.com/DavAnders/rss-aggregator/internal/database"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "ApiKey ") {
			RespondWithError(w, http.StatusBadRequest, "Invalid Authorization header")
			return
		}

		apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Unauthorized: Invalid API Key")
			return
		}

		handler(w, r, user)
	}
}
