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
	job, err := h.GetUC.Execute(id)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(job)
}
