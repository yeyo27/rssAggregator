package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/yeyo27/rssAggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string `json:"name"`
}


// this way we will return a normal user with the fields names specified in the json tags, using snake_case
func databaseUserToUser(dbUser database.User) User {
	return User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.CreatedAt,
		Name: dbUser.Name,
	}
}