package usecase

import (
	"context"
	"log"
	"os"
	"time"

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

func (u *ProcessJobUseCase) Execute(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	job, err := u.Queue.Dequeue()
	if err != nil {
		return err
	}

	log.Printf("PROCESSING JOB: %s (mime: %s)", job.ID, job.MimeType)

	defer func() {
		if job.Status == entity.JobDone || job.Status == entity.JobFailed {
			if err := os.Remove(job.InputPath); err != nil {
				log.Printf("Failed to cleanup input file %s: %v", job.InputPath, err)
			}

			if job.Status == entity.JobDone {
				time.AfterFunc(24*time.Hour, func() {
					if err := os.Remove(job.OutputPath); err != nil {
						log.Printf("Failed to cleanup output file %s: %v", job.OutputPath, err)
					}
				})
			}
		}
	}()

	job.Status = entity.JobRunning
	job.Progress = 10
	if err := u.Repository.Update(job); err != nil {
		log.Printf("Failed to update job status: %v", err)
	}

	select {
	case <-ctx.Done():
		job.Status = entity.JobFailed
		job.Error = "job cancelled"
		u.Repository.Update(job)
		return ctx.Err()
	default:
	}

	optimizer, err := u.FactoryFn(job.MimeType)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		if err := u.Repository.Update(job); err != nil {
			log.Printf("Failed to update job status: %v", err)
		}
		return err
	}

	job.Progress = 50
	if err := u.Repository.Update(job); err != nil {
		log.Printf("Failed to update job progress: %v", err)
	}

	select {
	case <-ctx.Done():
		job.Status = entity.JobFailed
		job.Error = "job cancelled"
		u.Repository.Update(job)
		return ctx.Err()
	default:
	}

	err = optimizer.Optimize(job.InputPath, job.OutputPath)
	if err != nil {
		job.Status = entity.JobFailed
		job.Error = err.Error()
		if err := u.Repository.Update(job); err != nil {
			log.Printf("Failed to update job status: %v", err)
		}
		return err
	}

	job.Status = entity.JobDone
	job.Progress = 100
	if err := u.Repository.Update(job); err != nil {
		log.Printf("Failed to update job status: %v", err)
	}

	log.Printf("JOB DONE: %s", job.ID)
	return nil
}
