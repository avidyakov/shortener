package handlers

import (
	"fmt"
	"github.com/avidyakov/shortener/internal/config"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/avidyakov/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

var Repo repositories.LinkRepo

func CreateShortLink(res http.ResponseWriter, req *http.Request) {
	originLink, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLinkID := utils.GenerateShortID(8)
	Repo.CreateLink(shortLinkID, string(originLink))

	shortLink := fmt.Sprintf("%s/%s", config.Cfg.BaseURL, shortLinkID)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))

	log.Printf("Short link created: %s -> %s", shortLink, originLink)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := chi.URLParam(r, "slug")
	originLink, ok := Repo.GetLink(shortLinkID)

	if ok {
		http.Redirect(w, r, originLink, http.StatusTemporaryRedirect)
		log.Printf("Redirected: %s/%s -> %s", config.Cfg.BaseURL, shortLinkID, originLink)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Short link %s/%s not found", config.Cfg.BaseURL, shortLinkID)
	}
}

func LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", CreateShortLink)
	r.Get("/{slug}", Redirect)
	return r
}
