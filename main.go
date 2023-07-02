package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/yeyo27/rssAggregator/internal/database"
	_ "github.com/lib/pq"
)

/*
Chi routers are used to define endpoints of the api.

Mount concatenates endpoints.

HandleFunc indicates what function is executed by the given endpoint. However,
that function is executed no matter what http method is used.
*/

type apiConfig struct {
	DB *database.Queries // connection to db
}

func main() {
	godotenv.Load() // ".env" by default

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not defined in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not defined in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:" ,err)
	}
	
	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness) // /v1/ready
	v1Router.Get("/error", handlerError) // /v1/error
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + portString,

	}

	log.Printf("Server starting on port %v\n", portString)
	err = server.ListenAndServe() // blocking function
	if err != nil {
		log.Fatal(err)
	}
}