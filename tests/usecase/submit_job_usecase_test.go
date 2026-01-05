package usecase_test

import (
	"testing"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestSubmitJobUseCase(t *testing.T) {
	queue := &mocks.JobQueueMock{}
	uc := usecase.NewSubmitJobUseCase(queue)

	job := entity.Job{
		ID: "job-1",
	}

	err := uc.Execute(job)
	if err != nil {
		t.Fatal("unexpected error")
	}

	if queue.StoredJob.Status != entity.JobPending {
		t.Errorf("expected PENDING, got %s", queue.StoredJob.Status)
	}
}
