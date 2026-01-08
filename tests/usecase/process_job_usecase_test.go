package usecase_test

import (
	"testing"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestProcessJobSuccess(t *testing.T) {
	queue := &mocks.JobQueueMock{}

	factoryMock := func(mime string) (service.Optimizer, error) {
		return &mocks.OptimizerMock{}, nil
	}

	uc := usecase.NewProcessJobUseCase(queue, factoryMock)

	job := entity.Job{
		ID:         "job-1",
		InputPath:  "input.jpeg",
		OutputPath: "output.jpeg",
		MimeType:   "image/jpeg",
	}

	uc.Execute(job)

	if queue.StoredJob.Status != entity.JobDone {
		t.Errorf("expected DONE, got %s", queue.StoredJob.Status)
	}

	if queue.StoredJob.Progress != 100 {
		t.Errorf("expected progress 100, got %d", queue.StoredJob.Progress)
	}
}
