package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/DavAnders/rss-aggregator/internal/scraper"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("test_connection")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	batchSize := int32(10)
	log.Println("Successfully connected to database using connection string: ", connStr)
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()
	interval := 60 * time.Second
	// Start the worker in its own goroutine
	go scraper.Worker(ctx, interval, batchSize, queries)

	// Example: setup a simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "RSS Aggregator is running...")
	})
	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}
}
