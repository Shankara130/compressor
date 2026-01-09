package usecase_test

import (
	"testing"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/factory"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/Shankara130/compressor/internal/usecase/mocks"
)

func TestProcessJobFail(t *testing.T) {
	job := entity.Job{
		ID:       "job-2",
		MimeType: "unsupported/type",
	}

	queue := &mocks.JobQueueMock{
		Job: job,
	}

	repo := &mocks.JobRepositoryMock{}

	uc := usecase.NewProcessJobUseCase(queue, repo, factory.NewOptimizer)

	uc.Execute()

	if repo.UpdatedJob.Status != entity.JobFailed {
		t.Errorf("expected FAILED, got %s", repo.UpdatedJob.Status)
	}
}
