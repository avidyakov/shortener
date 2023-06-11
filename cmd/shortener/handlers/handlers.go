package handlers

import (
	"fmt"
	. "github.com/avidyakov/shortener/cmd/shortener/config"
	"github.com/avidyakov/shortener/cmd/shortener/repositories"
	"github.com/avidyakov/shortener/cmd/shortener/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

var repo = repositories.NewMemoryLink()

func CreateShortLink(res http.ResponseWriter, req *http.Request) {
	originLink, _ := io.ReadAll(req.Body)
	shortLinkID := utils.GenerateShortID(8)
	repo.CreateLink(shortLinkID, string(originLink))

	shortLink := fmt.Sprintf("%s/%s", Config.BaseURL, shortLinkID)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))

	log.Printf("Short link created: %s -> %s", shortLink, originLink)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := chi.URLParam(r, "slug")
	originLink, ok := repo.GetLink(shortLinkID)

	if ok {
		http.Redirect(w, r, originLink, http.StatusTemporaryRedirect)
		log.Printf("Redirected: %s/%s -> %s", Config.BaseURL, shortLinkID, originLink)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Short link %s/%s not found", Config.BaseURL, shortLinkID)
	}
}

func LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", CreateShortLink)
	r.Get("/{slug}", Redirect)
	return r
}
