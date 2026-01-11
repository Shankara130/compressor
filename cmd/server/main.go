package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shankara130/compressor/internal/config"
	httpdelivery "github.com/Shankara130/compressor/internal/delivery/http"
	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/infrastructure/repository"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.InputDir, 0755); err != nil {
		log.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

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

	submitUC := usecase.NewSubmitJobUseCase(jobQueue, jobRepo)
	getUC := usecase.NewGetJobUseCase(jobRepo)

	uploadHandler := &handler.UploadHandler{SubmitUC: submitUC}
	statusHandler := &handler.StatusHandler{GetUC: getUC}
	downloadHandler := &handler.DownloadHandler{GetUC: getUC}

	router := httpdelivery.NewRouter(uploadHandler, statusHandler, downloadHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("UI server running at :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exited")
}
