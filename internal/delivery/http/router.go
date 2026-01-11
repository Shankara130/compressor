package http

import (
	"net/http"

	"github.com/Shankara130/compressor/internal/delivery/http/handler"
)

func NewRouter(
	uploadHandler *handler.UploadHandler,
	statusHandler *handler.StatusHandler,
	downloadHandler *handler.DownloadHandler,
	healthHandler *handler.HealthHandler,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", uploadHandler.Index)
	mux.HandleFunc("/upload", uploadHandler.Upload)
	mux.HandleFunc("/status/", statusHandler.Get)
	mux.HandleFunc("/download/", downloadHandler.Get)
	mux.HandleFunc("/health", healthHandler.Get)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
