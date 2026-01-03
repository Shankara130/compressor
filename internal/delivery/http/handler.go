package http

import (
	"net/http"

	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/usecase"
	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	file, _ := c.FormFile("file")
	input := "/tmp/input_" + file.Filename
	output := "/tmp/output_" + file.Filename

	c.SaveUploadedFile(file, input)

	uc := usecase.NewOptimizeFileUseCase()
	err := uc.Execute(entity.File{
		InputPath:  input,
		OutputPath: output,
		MimeType:   file.Header.Get("Content-Type"),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.File(output)
}
