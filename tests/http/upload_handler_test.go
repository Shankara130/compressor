package http_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestUploadHandler(t *testing.T) {
	queue := &mocks.JobQueueMock{}
	submitUC := usecase.NewSubmitJobUseCase(queue)

	h := &handler.UploadHandler{SubmitUC: submitUC}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("fake image"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	h.Upload(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
