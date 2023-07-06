package main

import (
	"fmt"
	"net/http"

	"github.com/yeyo27/rssAggregator/internal/database"
	"github.com/yeyo27/rssAggregator/internal/database/auth"
)

// this function type does not match the signature of an HTTP handler function,
// which is only made of http response writer and http request.
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// so this function must return an http.HandlerFunc to be used by the chi router
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't find user: %v", err))
			return
		}

		handler(w, r, user)
	}
}