package main

import (
	"log"
	"time"

	"github.com/Shankara130/compressor/internal/domain/factory"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/infrastructure/repository"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	jobQueue := queue.NewRedisQueue(redisClient)
	jobRepo := repository.NewRedisJobRepository(redisClient)

	processUC := usecase.NewProcessJobUseCase(jobQueue, jobRepo, factory.NewOptimizer)

	log.Println("worker started")

	for {
		err := processUC.Execute()
		if err != nil {
			log.Println("no job / error:", err)
			time.Sleep(time.Second)
		}
	}

}
