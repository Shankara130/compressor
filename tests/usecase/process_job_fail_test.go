package usecase_test

import (
	"testing"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestProcessJobFail(t *testing.T) {
	queue := &mocks.JobQueueMock{}
	uc := usecase.NewProcessJobUseCase(queue)

	job := entity.Job{
		ID:       "job-2",
		MimeType: "unsupported/type",
	}

	uc.Execute(job)

	if queue.StoredJob.Status != entity.JobFailed {
		t.Errorf("expected FAILED, got %s", queue.StoredJob.Status)
	}
}
