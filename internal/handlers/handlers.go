package handlers

import (
	"github.com/avidyakov/shortener/internal/encoding"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	baseURL string
	repo    repositories.LinkRepo
}

func NewHandlers(repo repositories.LinkRepo, baseURL string) *Handlers {
	return &Handlers{
		baseURL: baseURL,
		repo:    repo,
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
	r.Get("/ping", h.Ping)
	return r
}
