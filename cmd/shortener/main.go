package main

import (
	"github.com/avidyakov/shortener/internal/config"
	"github.com/avidyakov/shortener/internal/handlers"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger.Log, _ = zap.NewProduction()
	defer logger.Log.Sync()

	logger.Log.Info("Initializing server configuration and handlers")
	conf := config.NewConfig()
	handler := handlers.NewLinkHandlers(repositories.NewFileRepo(conf.File), conf.BaseURL)

	logger.Log.Info("Starting server",
		zap.String("serverAddr", conf.ServerAddr),
	)
	http.ListenAndServe(conf.ServerAddr, handler.LinkRouter())
}
