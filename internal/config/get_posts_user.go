package config

import (
	"net/http"
	"strconv"

	"github.com/DavAnders/rss-aggregator/internal/database"
)

func (cfg *ApiConfig) GetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// Default limit
	limit := 10

	// Check if a limit query parameter is provided
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	// Fetch posts
	posts, err := cfg.DB.GetPostsByUsers(r.Context(), database.GetPostsByUsersParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	// Respond with posts
	RespondWithJSON(w, http.StatusOK, posts)
}
