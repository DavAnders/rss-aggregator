package config

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB *database.Queries
}

func (cfg *ApiConfig) GetUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}
	RespondWithJSON(w, http.StatusOK, user)
}

func (cfg *ApiConfig) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
			return
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Error reading request body")
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
			RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}

		RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
	}

}
