package health

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alyelalwany/github-tracker/internal/model"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := model.HealthResponse{Status: "ok", Version: "dev"}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to send ok response", "err", err)
	}
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ready"}); err != nil {
		slog.Error("failed to send ready response", "err", err)
	}
}
