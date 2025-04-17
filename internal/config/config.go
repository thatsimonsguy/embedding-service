package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ModelPath       string
	EmbeddingBinary string
	BatchSize       string
}

var AppConfig Config

func Load() {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("[INFO] No .env file found, using environment variables")
	}

	AppConfig = Config{
		ModelPath:       getEnv("MODEL_PATH", "./nomic-embed-text-v1.5.Q4_K_M.gguf"),
		EmbeddingBinary: getEnv("EMBEDDING_BINARY", "./build/bin/llama-embedding"),
		BatchSize:       getEnv("BATCH_SIZE", "4096"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[WARN] %s not set, defaulting to %s\n", key, fallback)
		return fallback
	}
	return val
}
