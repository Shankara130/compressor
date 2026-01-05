package http

import (
	"net/http"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Queue service.JobQueue
}

func NewHandler(q service.JobQueue) *Handler {
	return &Handler{Queue: q}
}

func (h *Handler) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	id := uuid.NewString()

	input := "tmp/input/" + id
	output := "tmp/output/" + id

	_ = c.SaveUploadedFile(file, input)

	job := entity.Job{
		ID:         id,
		InputPath:  input,
		OutputPath: output,
		MimeType:   file.Header.Get("Content-Type"),
	}

	uc := usecase.NewSubmitJobUseCase(h.Queue)
	_ = uc.Execute(job)

	c.JSON(http.StatusAccepted, gin.H{"job_id": id})
}

func (h *Handler) GetStatus(c *gin.Context) {
	job, err := h.Queue.Get(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, job)
}

func (h *Handler) Download(c *gin.Context) {
	job, err := h.Queue.Get(c.Param("id"))
	if err != nil || job.Status != entity.JobDone {
		c.JSON(404, gin.H{"error": "not ready"})
		return
	}
	c.File(job.OutputPath)
}
