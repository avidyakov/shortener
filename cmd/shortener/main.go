package main

import (
	. "github.com/avidyakov/shortener/cmd/shortener/config"
	"github.com/avidyakov/shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting server at %s", Config.ServerAddr)
	http.ListenAndServe(Config.ServerAddr, handlers.LinkRouter())
}
