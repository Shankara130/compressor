package service

import (
	"context"

	"github.com/Shankara130/compressor/internal/domain/entity"
)

type JobQueue interface {
	Enqueue(job entity.Job) error
	Dequeue(ctx context.Context) (entity.Job, error)
}
