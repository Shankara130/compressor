package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
)

type DownloadHandler struct {
	GetUC *usecase.GetJobUseCase
}

func (h *DownloadHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/download/")

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
	if err != nil || job.Status != entity.JobDone {
		http.NotFound(w, r)
		return
	}

	// Security: Ensure file is in expected directory
	absOutputPath, err := filepath.Abs(job.OutputPath)
	if err != nil {
		http.Error(w, "invalid file path", http.StatusInternalServerError)
		return
	}

	expectedDir, err := filepath.Abs("tmp/output")
	if err != nil {
		http.Error(w, "invalid output directory", http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(absOutputPath, expectedDir) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, job.OutputPath)
}
