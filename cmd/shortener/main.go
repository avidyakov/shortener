package main

import (
	"github.com/avidyakov/shortener/cmd/shortener/config"
	"github.com/avidyakov/shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	log.Printf("Initializing config")
	config.Cfg = config.NewConfig()
	log.Printf("Starting server at %s", config.Cfg.ServerAddr)
	http.ListenAndServe(config.Cfg.ServerAddr, handlers.LinkRouter())
}
