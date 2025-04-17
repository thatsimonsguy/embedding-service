package main

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"matthewpsimons.com/embedding-service/handlers"
	"matthewpsimons.com/embedding-service/internal/config"
	"matthewpsimons.com/embedding-service/internal/logging"
)

func main() {
	config.Load()
	logger := logging.Init()
	defer logger.Sync()

	batchSize := config.AppConfig.BatchSize
	logger.Info("Parsed config",
		zap.String("model_path", config.AppConfig.ModelPath),
		zap.String("embedding_binary", config.AppConfig.EmbeddingBinary),
		zap.Int("batch_size", batchSize),
	)

	cfg := config.Config{
		ModelPath:       config.AppConfig.ModelPath,
		EmbeddingBinary: config.AppConfig.EmbeddingBinary,
		BatchSize:       strconv.Itoa(batchSize),
	}

	http.HandleFunc("/api/v1/embed", handlers.HandleEmbed(logger, cfg))

	logger.Info("Server starting", zap.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
