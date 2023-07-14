package handlers

import (
	"github.com/avidyakov/shortener/internal/encoding"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	baseURL     string
	databaseDSN string
	repo        repositories.LinkRepo
}

func NewHandlers(repo repositories.LinkRepo, baseURL, databaseDSN string) *Handlers {
	return &Handlers{
		baseURL:     baseURL,
		databaseDSN: databaseDSN,
		repo:        repo,
	}
}

func (h *Handlers) LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(encoding.GZIPMiddleware)
	r.Post("/", h.CreateShortLink)
	r.Post("/api/shorten", h.CreateShortLink) // for tests only
	r.Get("/{slug}", h.Redirect)
	r.Get("/ping", h.Ping)
	return r
}
