package queue

import (
	"context"
	"encoding/json"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	Client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{Client: client}
}

func (q *RedisQueue) Enqueue(job entity.Job) error {
	data, _ := json.Marshal(job)
	return q.Client.RPush(context.Background(), "jobs", data).Err()
}

func (q *RedisQueue) Dequeue() (entity.Job, error) {
	res, err := q.Client.BLPop(context.Background(), 0, "jobs").Result()
	if err != nil {
		return entity.Job{}, err
	}

	var job entity.Job
	_ = json.Unmarshal([]byte(res[1]), &job)
	return job, nil
}

func (q *RedisQueue) Update(job entity.Job) error {
	data, _ := json.Marshal(job)
	return q.Client.Set(context.Background(), "job:"+job.ID, data, 0).Err()
}

func (q *RedisQueue) Get(id string) (entity.Job, error) {
	val, err := q.Client.Get(context.Background(), "job:"+id).Result()
	if err != nil {
		return entity.Job{}, err
	}

	var job entity.Job
	_ = json.Unmarshal([]byte(val), &job)
	return job, nil
}
