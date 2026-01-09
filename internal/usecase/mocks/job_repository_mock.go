package mocks

import (
	"errors"

	"github.com/Shankara130/compressor/internal/domain/entity"
)

type JobRepositoryMock struct {
	UpdatedJob entity.Job
}

func (m *JobRepositoryMock) Save(job entity.Job) error {
	m.UpdatedJob = job
	return nil
}

func (m *JobRepositoryMock) Update(job entity.Job) error {
	m.UpdatedJob = job
	return nil
}

func (m *JobRepositoryMock) GetByID(id string) (entity.Job, error) {
	if m.UpdatedJob.ID != id {
		return entity.Job{}, errors.New("not found")
	}
	return m.UpdatedJob, nil
}
