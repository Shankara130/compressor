package http_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestUploadHandler(t *testing.T) {
	queue := &mocks.JobQueueMock{}
	submitUC := usecase.NewSubmitJobUseCase(queue, &mocks.JobRepositoryMock{})

	h := &handler.UploadHandler{SubmitUC: submitUC}

	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write(jpeg)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	h.Upload(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", rr.Code, rr.Body.String())
	}

	// Regression guard: the saved file on disk must be byte-identical to the
	// upload. Catches the streaming handler dropping the sniffed prefix.
	var resp struct {
		JobID string `json:"job_id"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parse response: %v (body: %s)", err, rr.Body.String())
	}
	saved, err := os.ReadFile("tmp/input/" + resp.JobID + ".jpg")
	if err != nil {
		t.Fatalf("read saved input file: %v", err)
	}
	if !bytes.Equal(saved, jpeg) {
		t.Fatalf("saved file (%d bytes) does not match upload (%d bytes)",
			len(saved), len(jpeg))
	}
}

func TestUploadHandler_RejectsUnsupportedType(t *testing.T) {
	queue := &mocks.JobQueueMock{}
	submitUC := usecase.NewSubmitJobUseCase(queue, &mocks.JobRepositoryMock{})
	h := &handler.UploadHandler{SubmitUC: submitUC}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "note.txt")
	part.Write([]byte("just some plain text"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	h.Upload(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unsupported type, got %d", rr.Code)
	}
}
