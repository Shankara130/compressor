package usecase_test

import (
	"testing"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestProcessJobSuccess(t *testing.T) {
	job := entity.Job{
		ID:         "job-1",
		InputPath:  "input.jpeg",
		OutputPath: "output.jpeg",
		MimeType:   "image/jpeg",
	}

	queue := &mocks.JobQueueMock{
		Job: job,
	}

	repo := &mocks.JobRepositoryMock{}

	factoryMock := func(mime string) (service.Optimizer, error) {
		return &mocks.OptimizerMock{}, nil
	}

	uc := usecase.NewProcessJobUseCase(queue, repo, factoryMock)

	uc.Execute()

	if repo.UpdatedJob.Status != entity.JobDone {
		t.Errorf("expected DONE, got %s", repo.UpdatedJob.Status)
	}

	if repo.UpdatedJob.Progress != 100 {
		t.Errorf("expected progress 100, got %d", repo.UpdatedJob.Progress)
	}
}
