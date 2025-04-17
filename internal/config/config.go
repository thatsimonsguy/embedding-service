package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ModelPath       string
	EmbeddingBinary string
	BatchSize       string
}

var AppConfig struct {
	ModelPath       string
	EmbeddingBinary string
	BatchSize       int
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[INFO] No .env file found, using environment variables")
	}

	AppConfig.ModelPath = getEnv("MODEL_PATH", "./nomic-embed-text-v1.5.Q4_K_M.gguf")
	AppConfig.EmbeddingBinary = getEnv("EMBEDDING_BINARY", "./build/bin/llama-embedding")
	AppConfig.BatchSize = getEnvInt("BATCH_SIZE", 4096)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[WARN] %s not set, defaulting to %s\n", key, fallback)
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		log.Printf("[WARN] %s not set, defaulting to %d\n", key, fallback)
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("[WARN] %s could not be parsed as int, defaulting to %d\n", key, fallback)
		return fallback
	}
	return val
}
