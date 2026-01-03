package factory

import (
	"errors"
	"strings"

	"github.com/Shankara130/compressor/internal/domain/service"
	infra "github.com/Shankara130/compressor/internal/infrastructure/optimizer"
)

func NewOptimizer(mime string) (service.Optimizer, error) {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return infra.NewImageOptimizer(), nil
	case mime == "application/pdf":
		return infra.NewPDFOptimizer(), nil
	case strings.HasPrefix(mime, "video/"):
		return infra.NewVideoOptimizer(), nil
	default:
		return nil, errors.New("unsupported file type")
	}
}
