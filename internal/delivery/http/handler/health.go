package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	RedisClient *redis.Client
}

type HealthResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	response := HealthResponse{
		Status: "healthy",
		Checks: make(map[string]string),
	}

	if err := h.RedisClient.Ping(ctx).Err(); err != nil {
		response.Status = "unhealthy"
		response.Checks["redis"] = "down: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Checks["redis"] = "up"
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
