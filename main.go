package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default to 8080
	}
	port = ":" + port

	mux := http.NewServeMux()

	wrappedMux := corsMiddleware(http.HandlerFunc(mainHandler))

	mux.Handle("/", wrappedMux)

	log.Fatal(http.ListenAndServe(port, mux))
}
