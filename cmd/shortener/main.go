package main

import (
	"github.com/avidyakov/shortener/cmd/shortener/config"
	"github.com/avidyakov/shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	log.Printf("Initializing config")
	config.Config = config.NewConfig()
	log.Printf("Starting server at %s", config.Config.ServerAddr)
	http.ListenAndServe(config.Config.ServerAddr, handlers.LinkRouter())
}
