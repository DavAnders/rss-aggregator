package config

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
)

func (cfg *ApiConfig) CreateFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	var req struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdAt := time.Now()
	updatedAt := createdAt

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      req.Name,
		Url:       req.Url,
		UserID:    user.ID,
	})

	if err != nil {
		log.Printf("Failed to create feed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	RespondWithJSON(w, http.StatusOK, feed)
}
