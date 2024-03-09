package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/DavAnders/rss-aggregator/internal/config"
	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	router := chi.NewRouter()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default to 8080
	}

	dbURL := os.Getenv("connection")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't open database: %v", err)
	}
	dbQueries := database.New(db)

	cfg := &config.ApiConfig{
		DB: dbQueries,
	}

	router.Use(config.CorsMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		config.MainHandler(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()
	v1Router.Post("/users", cfg.CreateUserHandler())
	v1Router.HandleFunc("/readiness", config.HandlerReadiness)
	v1Router.HandleFunc("/err", config.HandlerErr)

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", port)
	log.Fatal(srv.ListenAndServe())

}
