package main

import (
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	q := queue.NewRedisQueue(redisClient)
	uc := usecase.NewProcessJobUseCase(q)

	for {
		job, err := q.Dequeue()
		if err != nil {
			continue
		}
		uc.Execute(job)
	}
}
