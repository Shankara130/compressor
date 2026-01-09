package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/google/uuid"
)

type UploadHandler struct {
	SubmitUC *usecase.SubmitJobUseCase
}

func (h *UploadHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, nil)
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.New().String()
	input := "tmp/input/" + id
	output := "tmp/output/" + id

	os.MkdirAll("tmp/input", 0755)
	os.MkdirAll("tmp/output", 0755)

	dst, _ := os.Create(input)

	_, _ = io.Copy(dst, file)

	job := entity.Job{
		ID:         id,
		InputPath:  input,
		OutputPath: output,
		MimeType:   header.Header.Get("Content-Type"),
	}

	if err := h.SubmitUC.Execute(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"job_id": id,
	})
}
