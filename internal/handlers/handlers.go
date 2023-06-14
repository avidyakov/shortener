package handlers

import (
	"fmt"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/avidyakov/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type LinkHandlers struct {
	repo    repositories.LinkRepo
	baseURL string
}

func NewLinkHandlers(repo repositories.LinkRepo, baseURL string) *LinkHandlers {
	return &LinkHandlers{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (h *LinkHandlers) CreateShortLink(res http.ResponseWriter, req *http.Request) {
	originLink, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLinkID := utils.GenerateShortID(8)
	h.repo.CreateLink(shortLinkID, string(originLink))

	shortLink := fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))

	log.Printf("Short link created: %s -> %s", shortLink, originLink)
}

func (h *LinkHandlers) Redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := chi.URLParam(r, "slug")
	originLink, ok := h.repo.GetLink(shortLinkID)

	if ok {
		http.Redirect(w, r, originLink, http.StatusTemporaryRedirect)
		log.Printf("Redirected: %s/%s -> %s", h.baseURL, shortLinkID, originLink)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Short link %s/%s not found", h.baseURL, shortLinkID)
	}
}

func (h *LinkHandlers) LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateShortLink)
	r.Get("/{slug}", h.Redirect)
	return r
}
