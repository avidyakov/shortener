package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
)

var links = make(map[string]string)
var baseURL = "http://localhost:8080"

func generateShortID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)

	for i, value := range randomBytes {
		randomBytes[i] = charset[value%byte(len(charset))]
	}
	return string(randomBytes)
}

func createShortLink(res http.ResponseWriter, req *http.Request) {
	originLink, _ := io.ReadAll(req.Body)
	shortLinkID := generateShortID(8)
	links[shortLinkID] = string(originLink)
	shortLink := fmt.Sprintf("%s/%s", baseURL, shortLinkID)
	res.Write([]byte(shortLink))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := r.URL.Path[1:]
	originLink, ok := links[shortLinkID]

	if ok {
		http.Redirect(w, r, originLink, http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createShortLink(w, r)
	case http.MethodGet:
		redirect(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handle)
	http.ListenAndServe(`:8080`, mux)
}
