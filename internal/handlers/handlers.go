package handlers

import (
	"github.com/avidyakov/shortener/internal/encoding"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	baseURL   string
	repo      repositories.LinkRepo
	secretKey string
}

func NewHandlers(repo repositories.LinkRepo, baseURL, secretKey string) *Handlers {
	return &Handlers{
		baseURL:   baseURL,
		repo:      repo,
		secretKey: secretKey,
	}
}

func (h *Handlers) LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(encoding.GZIPMiddleware)
	r.Get("/{slug}", h.Redirect)
	r.Post("/", h.CreateShortLink)
	r.Post("/api/shorten", h.CreateShortLink) // for tests only
	r.Post("/api/shorten/batch", h.CreateABunchOfLinks)
	r.Get("/api/user/urls", h.UserURLs)
	r.Delete("/api/user/urls", h.DeleteUserURLs)
	r.Get("/ping", h.Ping)
	return r
}
