package handler

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/google/uuid"
)

type UploadHandler struct {
	SubmitUC *usecase.SubmitJobUseCase
}

func (h *UploadHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func outputExtension(mime string) string {
	switch {
	case strings.HasPrefix(mime, "video/"):
		return ".mp4"
	case mime == "image/jpeg":
		return ".jpg"
	case mime == "image/png":
		return ".png"
	case mime == "application/pdf":
		return ".pdf"
	default:
		return ""
	}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(500 << 20)
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
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "failed to reset file position", http.StatusInternalServerError)
		return
	}

	mime := http.DetectContentType(buf[:n])

	if err := ValidateFile(mime, header.Size, header.Filename); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, ErrFileTooLarge) {
			statusCode = http.StatusRequestEntityTooLarge
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	ext := outputExtension(mime)
	if ext == "" {
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	inputExt := strings.ToLower(filepath.Ext(header.Filename))
	input := "tmp/input/" + id + inputExt
	output := "tmp/output/" + id + ext

	if err := os.MkdirAll("tmp/input", 0755); err != nil {
		http.Error(w, "failed to create input directory", http.StatusInternalServerError)
		log.Printf("MkdirAll error: %v", err)
		return
	}

	if err := os.MkdirAll("tmp/output", 0755); err != nil {
		http.Error(w, "failed to create output directory", http.StatusInternalServerError)
		log.Printf("MkdirAll error: %v", err)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"job_id": id,
	}); err != nil {
		log.Printf("JSON encode error: %v", err)
	}

	log.Println("UPLOAD RECEIVED:", header.Filename)
}
