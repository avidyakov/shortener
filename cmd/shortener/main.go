package main

import (
	"github.com/avidyakov/shortener/internal/config"
	"github.com/avidyakov/shortener/internal/handlers"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"go.uber.org/zap"
	"net/http"
)

func getRepository(conf *config.Configuration) repositories.LinkRepo {
	if conf.DatabaseDSN != "" {
		return repositories.NewDBRepo(conf.DatabaseDSN)
	}

	if conf.File != "" {
		return repositories.NewFileRepo(conf.File)
	}

	return repositories.NewMemoryRepo()
}

func main() {
	logger.Log, _ = zap.NewProduction()
	defer logger.Log.Sync()

	logger.Log.Info("Initializing server configuration and handlers")
	conf := config.NewConfig()
	repo := getRepository(conf)
	handler := handlers.NewHandlers(repo, conf.BaseURL)

	logger.Log.Info("Starting server",
		zap.String("serverAddr", conf.ServerAddr),
	)
	http.ListenAndServe(conf.ServerAddr, handler.LinkRouter())
}
