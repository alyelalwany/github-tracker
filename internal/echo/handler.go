package echo

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/alyelalwany/github-tracker/internal/model"
)

/*
echoHandler only accepts POST requests. It parses the request and checks if it fits the echoRequest struct.
Upon success, it creates an echoResponse message, encodes it in JSON and sends it back as a response.
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) //Per RFC 7231
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB cap

	var req model.EchoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		slog.Warn("invalid request body", "err", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "field 'message' is required", http.StatusBadRequest)
		return
	}

	resp := model.EchoResponse{
		Message:    req.Message,
		ReceivedAt: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to write response", "err", err)
	}

}
