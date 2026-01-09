package service

import "github.com/Shankara130/compressor/internal/domain/entity"

type JobRepository interface {
	Save(job entity.Job) error
	Update(job entity.Job) error
	GetByID(id string) (entity.Job, error)
}
