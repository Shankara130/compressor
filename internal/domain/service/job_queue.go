package service

import "github.com/Shankara130/compressor/internal/domain/entity"

type JobQueue interface {
	Enqueue(job entity.Job) error
	Dequeue() (entity.Job, error)
}
