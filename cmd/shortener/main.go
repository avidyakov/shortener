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
	config.Cfg = config.NewConfig()

	log.Printf("Initializing handlers")
	handlers.Repo = repositories.NewMemoryLink()

	log.Printf("Starting server at %s", config.Cfg.ServerAddr)
	http.ListenAndServe(config.Cfg.ServerAddr, handlers.LinkRouter())
}
