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

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// input parameters to create user (the ones we must manually specify)
	type parameters struct {
		Name string `json:"name"`
	}

	// read params from request body
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 402, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// we want handlerGetUser to be an authenticated function, so it needs an user as an input parameter
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// we removed the logic to get the apikey and return the user by the given api key.
	// we are going to use a middleware function which implements that logic.
	// this way, we don't need to copy the auth part in every handler
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil {
		respondWithError(w, 402, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(posts))
}