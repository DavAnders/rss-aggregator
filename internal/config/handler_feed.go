package config

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/google/uuid"
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

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      req.Name,
		Url:       req.Url,
		UserID:    user.ID,
	})

	if err != nil {
		log.Printf("Failed to create feed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Failed to create feed follow: %v", err)
		feederr := cfg.DB.DeleteFeed(r.Context(), feed.ID)
		if feederr != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to delete feed")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to create feed follow")
		return
	}

	response := struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func (cfg *ApiConfig) GetFeedsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
			return
		}

		feeds, err := cfg.DB.GetAllFeeds(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch feed")
			return
		}
		RespondWithJSON(w, http.StatusOK, feeds)
	}
}
