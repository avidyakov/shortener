package main

import (
	"github.com/avidyakov/shortener/internal/config"
	"github.com/avidyakov/shortener/internal/handlers"
	"github.com/avidyakov/shortener/internal/repositories"
	"log"
	"net/http"
)

func main() {
	log.Printf("Initializing config")
	cfg := config.NewConfig()

	log.Printf("Initializing handlers")
	handler := handlers.NewLinkHandlers(repositories.NewMemoryLink(), cfg.BaseURL)

	log.Printf("Starting server at %s", cfg.ServerAddr)
	http.ListenAndServe(cfg.ServerAddr, handler.LinkRouter())
}
