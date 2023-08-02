package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/models"
	"github.com/avidyakov/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const TokenExp = time.Hour * 24

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func (h *Handlers) CreateShortLink(res http.ResponseWriter, req *http.Request) {
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
		originLink = model.URL
	default:
		originBytes, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Error("Invalid request body",
				zap.Error(err),
			)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		originLink = string(originBytes)
	}

	validatedLink, err := utils.ValidateLink(originLink)
	if err != nil {
		logger.Log.Error("Invalid link",
			zap.String("originLink", originLink),
		)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortLinkID := utils.GenerateShortID(8)
	// get cookie
	parsedToken, _ := req.Cookie("token")
	userID := -1
	if parsedToken != nil {
		userID = h.getUserID(parsedToken.Value)
	}

	if userID == -1 {
		userID, err = h.repo.CreateUser()
	}

	if err != nil {
		logger.Log.Error("Can't create user",
			zap.Error(err),
		)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.repo.CreateLink(shortLinkID, validatedLink, userID)
	status := http.StatusCreated
	if err != nil {
		status = http.StatusConflict
		shortLinkID, _ = h.repo.GetShortLink(validatedLink)
	}

	shortLink := fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)
	responseData := shortLink

	if req.Header.Get("Content-Type") == "application/json" {
		res.Header().Set("Content-Type", "application/json")
		responseData = fmt.Sprintf(`{"result":"%s"}`, shortLink)
	}
	token, err := h.buildJWTString(userID)
	if err != nil {
		logger.Log.Error("Can't create JWT token",
			zap.Error(err),
		)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(TokenExp),
	})
	res.WriteHeader(status)
	res.Write([]byte(responseData))

	logger.Log.Info("Short link created",
		zap.String("shortLink", shortLink),
		zap.String("originLink", originLink),
	)
}

func (h *Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	shortLinkID := chi.URLParam(r, "slug")
	originLink, ok := h.repo.GetOriginLink(shortLinkID)

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

func (h *Handlers) CreateABunchOfLinks(w http.ResponseWriter, r *http.Request) {
	var longUrls []models.RequestLinkBatch
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&longUrls); err != nil {
		logger.Log.Debug("Invalid JSON",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var shortUrls []models.ResponseLinkBatch
	for _, link := range longUrls {
		validatedLink, err := utils.ValidateLink(link.OriginURL)
		if err != nil {
			logger.Log.Error("Invalid link",
				zap.String("originLink", link.OriginURL),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortLinkID := utils.GenerateShortID(8)
		userID := 1
		h.repo.CreateLink(shortLinkID, validatedLink, userID)

		shortLink := fmt.Sprintf("%s/%s", h.baseURL, shortLinkID)
		shortUrls = append(shortUrls, models.ResponseLinkBatch{
			CorrelationID: link.CorrelationID,
			ShortURL:      shortLink,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortUrls)
}

func (h *Handlers) buildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	return token.SignedString([]byte(h.secretKey))
}
