package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"mime/multipart"
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

// detectMime returns the content type of the uploaded file. It prefers the
// sniffed type, but falls back to the extension when the content is opaque —
// Go's sniffer cannot identify some video containers such as QuickTime (.mov),
// which would otherwise be rejected as unsupported even though they're allowed.
func detectMime(sniff []byte, filename string) string {
	if m := http.DetectContentType(sniff); m != "application/octet-stream" {
		return m
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".mp4", ".m4v":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".webm":
		return "video/webm"
	case ".pdf":
		return "application/pdf"
	}
	return "application/octet-stream"
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Cap the total request body at the largest allowed upload size, so
	// oversized requests are rejected early instead of being buffered.
	r.Body = http.MaxBytesReader(w, r.Body, MaxVideoSize)

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	// Walk the parts until we find the "file" field.
	var part *multipart.Part
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			http.Error(w, "file field is required", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "failed to read upload", http.StatusBadRequest)
			return
		}
		if p.FormName() == "file" {
			part = p
			break
		}
	}
	defer part.Close()

	// Sniff the content type from the first 512 bytes.
	sniff := make([]byte, 512)
	n, err := io.ReadFull(part, sniff)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	mime := detectMime(sniff[:n], part.FileName())

	// Validate type and extension up front; the per-type size limit is
	// enforced while streaming below.
	if err := ValidateFile(mime, 0, part.FileName()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ext := outputExtension(mime)
	if ext == "" {
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	inputExt := strings.ToLower(filepath.Ext(part.FileName()))
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

	// Stream the file to disk: write the sniffed prefix first (those bytes
	// were already consumed from the part), then copy the rest while enforcing
	// the per-type size limit. Oversized uploads are rejected as they arrive
	// rather than after the whole file is written.
	maxSize := MaxSizeForMime(mime)
	var copied int64
	var copyErr error
	if _, copyErr = dst.Write(sniff[:n]); copyErr == nil {
		copied, copyErr = io.Copy(dst, io.LimitReader(part, maxSize-int64(n)+1))
	}
	if cerr := dst.Close(); cerr != nil && copyErr == nil {
		copyErr = cerr
	}
	if copyErr != nil {
		os.Remove(input)
		http.Error(w, "failed to store file", http.StatusInternalServerError)
		return
	}
	if total := int64(n) + copied; total > maxSize {
		os.Remove(input)
		http.Error(w, "file too large: exceeded size limit for this file type", http.StatusRequestEntityTooLarge)
		return
	}

	job := entity.Job{
		ID:         id,
		InputPath:  input,
		OutputPath: output,
		MimeType:   mime,
	}

	if err := h.SubmitUC.Execute(job); err != nil {
		os.Remove(input)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"job_id": id,
	}); err != nil {
		log.Printf("JSON encode error: %v", err)
	}

	log.Println("UPLOAD RECEIVED:", part.FileName())
}
