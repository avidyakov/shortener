package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/models"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/avidyakov/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
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
	var originLink string
	switch req.Header.Get("Content-Type") {
	case "application/json":
		var model models.ShortURLRequest
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&model); err != nil {
			logger.Log.Debug("Invalid JSON",
				zap.Error(err),
			)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		originLink = model.Url
	case "text/plain":
		originBytes, _ := io.ReadAll(req.Body)
		originLink = string(originBytes)
	default:
		http.Error(res, "Unsupported content type", http.StatusUnsupportedMediaType)
	}

	validatedLink, err := utils.ValidateLink(originLink)
	if err != nil {
		logger.Log.Error("Invalid link",
			zap.String("originLink", string(originLink)),
		)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortLinkID := utils.GenerateShortID(8)
	h.repo.CreateLink(shortLinkID, validatedLink)

	shortLink := fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))

	logger.Log.Info("Short link created",
		zap.String("shortLink", shortLink),
		zap.String("originLink", string(originLink)),
	)
}

func (h *LinkHandlers) Redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := chi.URLParam(r, "slug")
	originLink, ok := h.repo.GetLink(shortLinkID)

	if ok {
		http.Redirect(w, r, originLink, http.StatusTemporaryRedirect)
		logger.Log.Info("Redirected",
			zap.String("shortLink", fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)),
			zap.String("originLink", originLink),
		)
	} else {
		w.WriteHeader(http.StatusNotFound)
		logger.Log.Info("Short link not found",
			zap.String("shortLink", fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)),
		)
	}
}

func (h *LinkHandlers) LinkRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Post("/", h.CreateShortLink)
	r.Get("/{slug}", h.Redirect)
	return r
}
