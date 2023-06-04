package main

import (
	"github.com/avidyakov/shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", handlers.LinkRouter())
}
