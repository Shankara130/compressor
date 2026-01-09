package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
)

type OptimizerFactoryFunc func(mime string) (service.Optimizer, error)

type ProcessJobUseCase struct {
	Queue      service.JobQueue
	Repository service.JobRepository
	FactoryFn  OptimizerFactoryFunc
}

func NewProcessJobUseCase(q service.JobQueue, r service.JobRepository, factoryFn OptimizerFactoryFunc) *ProcessJobUseCase {
	return &ProcessJobUseCase{Queue: q, Repository: r, FactoryFn: factoryFn}
}

func (u *ProcessJobUseCase) Execute() {
	job, err := u.Queue.Dequeue()
	if err != nil {
		return
	}

	job.Status = entity.JobRunning
	job.Progress = 10
	_ = u.Repository.Update(job)

	optimizer, err := u.FactoryFn(job.MimeType)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		_ = u.Repository.Update(job)
		return
	}

	job.Progress = 50
	_ = u.Repository.Update(job)

	err = optimizer.Optimize(job.InputPath, job.OutputPath)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		_ = u.Repository.Update(job)
		return
	}

	job.Status = entity.JobDone
	job.Progress = 100
	_ = u.Repository.Update(job)
}
