package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, _ = file.Read(buf)
	file.Seek(0, 0)

	mime := http.DetectContentType(buf)

	id := uuid.New().String()
	input := "tmp/input/" + id
	output := "tmp/output/" + id

	os.MkdirAll("tmp/input", 0755)
	os.MkdirAll("tmp/output", 0755)

	dst, err := os.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	job := entity.Job{
		ID:         id,
		InputPath:  input,
		OutputPath: output,
		MimeType:   mime,
	}

	if err := h.SubmitUC.Execute(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"job_id": id,
	})

	log.Println("UPLOAD RECEIVED:", header.Filename)
}
