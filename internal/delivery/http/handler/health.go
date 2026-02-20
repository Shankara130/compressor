package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "healthy",
		Checks: map[string]string{
			"queue":      "up (in-memory)",
			"repository": "up (in-memory)",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
