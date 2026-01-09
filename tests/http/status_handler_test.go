package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestStatusHandler(t *testing.T) {
	repo := &mocks.JobRepositoryMock{
		UpdatedJob: entity.Job{
			ID:     "123",
			Status: entity.JobDone,
		},
	}

	getUC := usecase.NewGetJobUseCase(repo)
	h := &handler.StatusHandler{GetUC: getUC}

	req := httptest.NewRequest("GET", "/status/123", nil)
	rr := httptest.NewRecorder()

	h.Get(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
