package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shankara130/compressor/internal/config"
	httpdelivery "github.com/Shankara130/compressor/internal/delivery/http"
	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/domain/factory"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/infrastructure/repository"
	"github.com/Shankara130/compressor/internal/usecase"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.InputDir, 0755); err != nil {
		log.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Shared queue and repository — used by BOTH the HTTP server and workers
	jobQueue := queue.NewInMemoryJobQueue()
	jobRepo := repository.NewInMemoryJobRepository()

	// Use cases
	submitUC := usecase.NewSubmitJobUseCase(jobQueue, jobRepo)
	getUC := usecase.NewGetJobUseCase(jobRepo)
	processUC := usecase.NewProcessJobUseCase(jobQueue, jobRepo, factory.NewOptimizer)

	// HTTP handlers
	uploadHandler := &handler.UploadHandler{SubmitUC: submitUC}
	statusHandler := &handler.StatusHandler{GetUC: getUC}
	downloadHandler := &handler.DownloadHandler{GetUC: getUC}
	healthHandler := &handler.HealthHandler{}

	router := httpdelivery.NewRouter(uploadHandler, statusHandler, downloadHandler, healthHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start workers in background goroutines
	var wg sync.WaitGroup
	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("Worker %d started", workerID)
			for {
				select {
				case <-ctx.Done():
					log.Printf("Worker %d shutting down", workerID)
					return
				default:
					if err := processUC.Execute(ctx); err != nil && ctx.Err() != nil {
						log.Printf("Worker %d error: %v", workerID, err)
					}
				}
			}
		}(i)
	}
	log.Printf("Started %d workers", cfg.WorkerCount)

	// Start HTTP server
	go func() {
		log.Printf("HTTP server running at :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Wait for workers to finish
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

	log.Println("Exited")
}
