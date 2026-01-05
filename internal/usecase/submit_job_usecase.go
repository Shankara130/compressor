package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
)

type SubmitJobUseCase struct {
	Queue service.JobQueue
}

func NewSubmitJobUseCase(q service.JobQueue) *SubmitJobUseCase {
	return &SubmitJobUseCase{Queue: q}
}

func (u *SubmitJobUseCase) Execute(job entity.Job) error {
	job.Status = entity.JobPending
	job.Progress = 0
	return u.Queue.Enqueue(job)
}
