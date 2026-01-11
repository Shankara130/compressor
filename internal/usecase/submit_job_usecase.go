package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
)

type SubmitJobUseCase struct {
	Queue      service.JobQueue
	Repository service.JobRepository
}

func NewSubmitJobUseCase(q service.JobQueue, r service.JobRepository) *SubmitJobUseCase {
	return &SubmitJobUseCase{Queue: q, Repository: r}
}

func (u *SubmitJobUseCase) Execute(job entity.Job) error {
	job.Status = entity.JobPending
	job.Progress = 0

	if err := u.Repository.Save(job); err != nil {
		return err
	}

	return u.Queue.Enqueue(job)
}
