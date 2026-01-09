package handler

import (
	"net/http"
	"strings"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
)

type DownloadHandler struct {
	GetUC *usecase.GetJobUseCase
}

func (h *DownloadHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/download/")
	job, err := h.GetUC.Execute(id)

	if err != nil || job.Status != entity.JobDone {
		w.WriteHeader(404)
		return
	}

	http.ServeFile(w, r, job.OutputPath)
}
