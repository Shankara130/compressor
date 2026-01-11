package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shankara130/compressor/internal/config"
	"github.com/Shankara130/compressor/internal/domain/factory"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/infrastructure/repository"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	jobQueue := queue.NewRedisQueue(redisClient)
	jobRepo := repository.NewRedisJobRepository(redisClient)

	processUC := usecase.NewProcessJobUseCase(jobQueue, jobRepo, factory.NewOptimizer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			log.Printf("Worker %d started", workerId)

			for {
				select {
				case <-ctx.Done():
					log.Printf("Worker %d shutting down", workerId)
					return
				default:
					if err := processUC.Execute(ctx); err != nil {
						if ctx.Err() != nil {
							log.Printf("Worker %d error: %v", workerID, err)
						}
					}
				}
			}
		}(i)
	}

	log.Printf("Started %d workers", cfg.WorkerCount)

	<-ctx.Done()
	log.Println("Shutting down workers gracefully...")

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Println("Workers shutdown timeout exceeded")
	}

	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Worker process exited")
}
