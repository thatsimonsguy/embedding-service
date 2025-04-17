package main

import (
	"net/http"

	"go.uber.org/zap"
	"matthewpsimons.com/embedding-service/handlers"
	"matthewpsimons.com/embedding-service/internal/config"
	"matthewpsimons.com/embedding-service/internal/logging"
)

func main() {
	config.Load()
	logging.Init()
	defer logging.Logger.Sync()

	log := logging.Logger

	cfg := config.Config{
		ModelPath:       config.AppConfig.ModelPath,
		EmbeddingBinary: config.AppConfig.EmbeddingBinary,
		BatchSize:       config.AppConfig.BatchSize,
	}

	log.Info("Loaded config",
		zap.String("model_path", cfg.ModelPath),
		zap.String("embedding_binary", cfg.EmbeddingBinary),
		zap.String("batch size", cfg.BatchSize),
	)

	http.HandleFunc("/api/v1/embed", handlers.HandleEmbed(log, cfg))

	log.Info("Server starting", zap.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
