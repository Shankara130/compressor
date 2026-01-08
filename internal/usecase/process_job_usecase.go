package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
)

type OptimizerFactoryFunc func(mime string) (service.Optimizer, error)

type ProcessJobUseCase struct {
	Queue     service.JobQueue
	FactoryFn OptimizerFactoryFunc
}

func NewProcessJobUseCase(q service.JobQueue, factoryFn OptimizerFactoryFunc) *ProcessJobUseCase {
	return &ProcessJobUseCase{Queue: q, FactoryFn: factoryFn}
}

func (u *ProcessJobUseCase) Execute(job entity.Job) {
	job.Status = entity.JobRunning
	job.Progress = 10
	_ = u.Queue.Update(job)

	optimizer, err := u.FactoryFn(job.MimeType)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		_ = u.Queue.Update(job)
		return
	}

	job.Progress = 50
	_ = u.Queue.Update(job)

	err = optimizer.Optimize(job.InputPath, job.OutputPath)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		_ = u.Queue.Update(job)
		return
	}

	job.Status = entity.JobDone
	job.Progress = 100
	_ = u.Queue.Update(job)
}
