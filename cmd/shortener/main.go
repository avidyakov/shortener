package main

import (
	"github.com/avidyakov/shortener/cmd/shortener/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.Handle)
	http.ListenAndServe(`:8080`, mux)
}
