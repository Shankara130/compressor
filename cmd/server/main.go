package main

import (
	"log"
	"net/http"
	"os"

	httpdelivery "github.com/Shankara130/compressor/internal/delivery/http"
	"github.com/Shankara130/compressor/internal/delivery/http/handler"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/Shankara130/compressor/internal/infrastructure/repository"
	"github.com/Shankara130/compressor/internal/usecase"
)

func main() {
	_ = os.MkdirAll("tmp/input", 0755)
	_ = os.MkdirAll("tmp/output", 0755)

	jobQueue := queue.NewInMemoryJobQueue()
	jobRepo := repository.NewInMemoryJobRepository()

	submitUC := usecase.NewSubmitJobUseCase(jobQueue)
	getUC := usecase.NewGetJobUseCase(jobRepo)

	uploadHandler := &handler.UploadHandler{SubmitUC: submitUC}
	statusHandler := &handler.StatusHandler{GetUC: getUC}
	downloadHandler := &handler.DownloadHandler{GetUC: getUC}

	router := httpdelivery.NewRouter(uploadHandler, statusHandler, downloadHandler)

	log.Println("UI server running at :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
