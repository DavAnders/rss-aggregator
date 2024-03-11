package config

import (
	"encoding/json"
	"net/http"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var req struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	db := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: req.FeedID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), db)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}
	RespondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *ApiConfig) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error fetching feed follows")
		return
	}
	apiFeedFollows := make([]FeedFollow, 0, len(feedFollows))
	for _, feedFollow := range feedFollows {
		apiFeedFollow := databaseFeedFollowToFeedFollow(feedFollow)
		apiFeedFollows = append(apiFeedFollows, apiFeedFollow)
	}

	RespondWithJSON(w, http.StatusOK, apiFeedFollows)
}

func (cfg *ApiConfig) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowID := chi.URLParam(r, "feedFollowID")

	// feedfollowid string to uuid
	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid feed follow ID format")
		return
	}

	feedFollow, err := cfg.DB.GetFeedFollowByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error retrieving feed follow by ID")
		return
	}

	if feedFollow.UserID != user.ID {
		RespondWithError(w, http.StatusForbidden, "You do not have permission to delete this feed follow")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
