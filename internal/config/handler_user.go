package config

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB *database.Queries
}

func (cfg *ApiConfig) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "ApiKey ") {
			RespondWithError(w, http.StatusBadRequest, "Invalid Authorization header")
			return
		}

		apiKey := strings.TrimPrefix(authHeader, "ApiKey ")

		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching user")
			return
		}

		RespondWithJSON(w, http.StatusOK, user)
	}
}

func (cfg *ApiConfig) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		id := uuid.New()
		createdAt := time.Now()
		updatedAt := time.Now()

		user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Name:      req.Name,
		})
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
		}
	}
}
