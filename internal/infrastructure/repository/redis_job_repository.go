package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

type RedisJobRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisJobRepository(client *redis.Client) *RedisJobRepository {
	if client == nil {
		panic("redis client cannot be nil")
	}

	return &RedisJobRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisJobRepository) key(id string) string {
	return fmt.Sprintf("job:%s", id)
}

func (r *RedisJobRepository) Save(job entity.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, r.key(job.ID), data, 0).Err()
}

func (r *RedisJobRepository) Update(job entity.Job) error {
	return r.Save(job)
}

func (r *RedisJobRepository) GetByID(id string) (entity.Job, error) {
	val, err := r.client.Get(r.ctx, r.key(id)).Result()
	if err == redis.Nil {
		return entity.Job{}, errors.New("job not found")
	}
	if err != nil {
		return entity.Job{}, err
	}

	var job entity.Job
	if err := json.Unmarshal([]byte(val), &job); err != nil {
		return entity.Job{}, err
	}

	return job, nil
}
