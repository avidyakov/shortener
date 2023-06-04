package handlers

import (
	"fmt"
	"github.com/avidyakov/shortener/cmd/shortener/repositories"
	"github.com/avidyakov/shortener/cmd/shortener/utils"
	"io"
	"net/http"
)

var baseURL = "http://localhost:8080"
var repo = repositories.NewMemoryLink()

func createShortLink(res http.ResponseWriter, req *http.Request) {
	originLink, _ := io.ReadAll(req.Body)
	shortLinkID := utils.GenerateShortID(8)
	repo.CreateLink(shortLinkID, string(originLink))

	shortLink := fmt.Sprintf("%s/%s", baseURL, shortLinkID)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := r.URL.Path[1:]
	originLink, ok := repo.GetLink(shortLinkID)

	if ok {
		http.Redirect(w, r, originLink, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createShortLink(w, r)
	case http.MethodGet:
		redirect(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
