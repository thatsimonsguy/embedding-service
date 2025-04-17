package handlers

import (
	"encoding/json"
	"fmt"
	"matthewpsimons.com/embedding-service/internal/config"
	"net/http"
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
		var req EmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		cmd := exec.Command(cfg.EmbeddingBinary,
			"-m", cfg.ModelPath,
			"--batch-size", cfg.BatchSize,
			"-p", req.Text,
		)
		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "embedding failed: "+string(outputBytes), http.StatusInternalServerError)
			return
		}

		embedding, err := extractEmbedding(string(outputBytes))
		if err != nil {
			http.Error(w, "parse failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(EmbedResponse{Embedding: embedding})
	}
}

func extractEmbedding(output string) ([]float32, error) {
	var result []float32
	start := strings.Index(output, "embedding 0:")
	if start == -1 {
		return nil, nil
	}
	line := output[start+len("embedding 0:"):]
	fields := strings.Fields(line)
	for _, field := range fields {
		if f, err := parseFloat(field); err == nil {
			result = append(result, f)
		}
	}
	return result, nil
}

func parseFloat(s string) (float32, error) {
	var f float32
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
