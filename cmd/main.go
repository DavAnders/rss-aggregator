package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/DavAnders/rss-aggregator/internal/config"
	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/DavAnders/rss-aggregator/internal/scraper"
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
	defer db.Close()
	dbQueries := database.New(db)

	cfg := &config.ApiConfig{
		DB: dbQueries,
	}

	router.Use(config.CorsMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		config.MainHandler(w, r)
	})

	// scraper
	ctx := context.Background()
	interval := 60 * time.Second // Customize as needed
	batchSize := int32(10)       // Customize as needed
	go scraper.Worker(ctx, interval, batchSize, dbQueries)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()
	v1Router.Post("/users", cfg.CreateUserHandler())
	v1Router.Get("/users", cfg.MiddlewareAuth(cfg.GetUserHandler))
	v1Router.HandleFunc("/readiness", config.HandlerReadiness)
	v1Router.HandleFunc("/err", config.HandlerErr)
	v1Router.Post("/feeds", cfg.MiddlewareAuth(cfg.CreateFeedHandler))
	v1Router.Get("/feeds", cfg.GetFeedsHandler())
	v1Router.Post("/feed_follows", cfg.MiddlewareAuth(cfg.CreateFeedFollowHandler))
	v1Router.Get("/feed_follows", cfg.MiddlewareAuth(cfg.GetFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", cfg.MiddlewareAuth(cfg.DeleteFeedFollowHandler))
	v1Router.Get("/posts", cfg.MiddlewareAuth(cfg.GetPostsByUser))

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", port)
	log.Fatal(srv.ListenAndServe())

}
