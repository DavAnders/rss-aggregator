package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DavAnders/rss-aggregator/internal/config"
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

	router.Use(config.CorsMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		config.MainHandler(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()

	v1Router.HandleFunc("/readiness", config.HandlerReadiness)
	v1Router.HandleFunc("/err", config.HandlerErr)

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", port)
	log.Fatal(srv.ListenAndServe())

}
