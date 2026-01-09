package mocks

import "github.com/Shankara130/compressor/internal/domain/entity"

type JobQueueMock struct {
	Job entity.Job
}

func (m *JobQueueMock) Enqueue(job entity.Job) error {
	m.Job = job
	return nil
}

func (m *JobQueueMock) Dequeue() (entity.Job, error) {
	return m.Job, nil
}

func (m *JobQueueMock) Update(job entity.Job) error {
	m.Job = job
	return nil
}

func (m *JobQueueMock) Get(id string) (entity.Job, error) {
	return m.Job, nil
}
