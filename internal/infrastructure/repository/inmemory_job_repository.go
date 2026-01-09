package repository

import (
	"errors"
	"sync"

	"github.com/Shankara130/compressor/internal/domain/entity"
)

type InMemoryJobRepository struct {
	mu   sync.Mutex
	jobs map[string]entity.Job
}

func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]entity.Job),
	}
}

func (r *InMemoryJobRepository) Save(job entity.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryJobRepository) Update(job entity.Job) error {
	return r.Save(job)
}

func (r *InMemoryJobRepository) GetByID(id string) (entity.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[id]
	if !ok {
		return entity.Job{}, errors.New("not found")
	}
	return job, nil
}
