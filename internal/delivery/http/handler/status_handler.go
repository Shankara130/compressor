package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Shankara130/compressor/internal/usecase"
)

type StatusHandler struct {
	GetUC *usecase.GetJobUseCase
}

func (h *StatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/status/")

	// Security: Prevent path traversal attacks
	if strings.Contains(id, "..") || strings.ContainsAny(id, "/") || strings.Contains(id, "\\") {
		http.Error(w, "invalid job ID", http.StatusBadRequest)
		return
	}

	if id == "" {
		http.NotFound(w, r)
		return
	}

	job, err := h.GetUC.Execute(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
