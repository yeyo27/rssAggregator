package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yeyo27/rssAggregator/internal/database"
)

// we need to grant acces to the database. The solution is turning this function into a method, owned by
// the apiConfig type, whose field DB contains the connection to the database

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// input parameters to create user (the ones we must manually specify)
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:	   params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 402, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

// we want handlerGetFeed to be an authenticated function, so it needs an user as an input parameter
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 402, fmt.Sprintf("Couldn't fetch feeds: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}
