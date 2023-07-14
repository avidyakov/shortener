package handlers

import (
	"context"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	conn, err := pgx.Connect(context.Background(), h.databaseDSN)
	var status int
	if err == nil {
		defer conn.Close(context.Background())
		status = http.StatusOK
		w.Write([]byte("pong"))
	} else {
		status = http.StatusInternalServerError
		logger.Log.Error("Ping failed", zap.Error(err))
	}
	w.WriteHeader(status)
}
