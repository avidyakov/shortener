package handlers

import (
	"github.com/avidyakov/shortener/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	err := h.repo.CheckConnection()
	var status int
	if err == nil {
		status = http.StatusOK
		w.Write([]byte("pong"))
	} else {
		status = http.StatusInternalServerError
		logger.Log.Error("Ping failed", zap.Error(err))
	}
	w.WriteHeader(status)
}
