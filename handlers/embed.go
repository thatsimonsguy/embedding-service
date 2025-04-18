package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"matthewpsimons.com/embedding-service/internal/config"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type EmbedRequest struct {
	Text string `json:"text"`
}

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func HandleEmbed(logger *zap.Logger, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req EmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			logger.Warn("Invalid request body", zap.Error(err))
			return
		}

		tmpfile, err := os.CreateTemp("", "prompt-*.txt")
		if err != nil {
			logger.Error("Failed to create temp file", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.WriteString(req.Text); err != nil {
			logger.Error("Failed to write prompt to temp file", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		tmpfile.Close()

		args := []string{"-m", cfg.ModelPath, "--batch-size", cfg.BatchSize, "--file", tmpfile.Name()}
		logger.Info("Running embedding command", zap.String("binary", cfg.EmbeddingBinary), zap.Strings("args", args))

		cmd := exec.Command(cfg.EmbeddingBinary, args...)
		outputBytes, err := cmd.CombinedOutput()
		output := string(outputBytes)
		logger.Info("Command output captured", zap.Int("output_length", len(output)), zap.String("tail", truncateEnd(output, 500)))

		if err != nil {
			http.Error(w, "embedding failed", http.StatusInternalServerError)
			logger.Error("embedding failed", zap.Error(err), zap.String("output_snippet", truncate(output, 500)))
			return
		}

		embedding, err := extractEmbedding(output)
		if err != nil {
			http.Error(w, "parse failed", http.StatusInternalServerError)
			logger.Error("failed to parse embedding", zap.Error(err), zap.String("raw_output", truncate(output, 1000)))
			return
		}

		logger.Info("Parsed embedding", zap.Int("dimensions", len(embedding)))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(EmbedResponse{Embedding: embedding})
	}
}

func extractEmbedding(output string) ([]float32, error) {
	var result []float32
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "embedding 0:") {
			fields := strings.Fields(line[len("embedding 0:"):])
			for _, field := range fields {
				if f, err := parseFloat(field); err == nil {
					result = append(result, f)
				}
			}
			break
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("embedding parse resulted in 0-length vector or marker not found")
	}
	return result, nil
}

func parseFloat(s string) (float32, error) {
	var f float32
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "... [truncated]"
	}
	return s
}

func truncateEnd(s string, max int) string {
	if len(s) > max {
		return "..." + s[len(s)-max:]
	}
	return s
}
