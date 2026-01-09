package queue

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
)

type InMemoryJobQueue struct {
	ch chan entity.Job
}

func NewInMemoryJobQueue() *InMemoryJobQueue {
	return &InMemoryJobQueue{
		ch: make(chan entity.Job, 100),
	}
}

func (q *InMemoryJobQueue) Enqueue(job entity.Job) error {
	q.ch <- job
	return nil
}

func (q *InMemoryJobQueue) Dequeue() (entity.Job, error) {
	return <-q.ch, nil
}
