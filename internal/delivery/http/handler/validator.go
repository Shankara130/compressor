package handler

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	MaxImageSize = 50 << 20  // 50 MB
	MaxVideoSize = 500 << 20 // 500 MB
	MaxPDFSize   = 50 << 20  // 50 MB
)

var (
	ErrFileTooLarge      = errors.New("file too large")
	ErrUnsupportedFormat = errors.New("unsupported file format")
	ErrUnsupportedType   = errors.New("unsupported file type")
)

func ValidateFile(mime string, size int64, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	switch {
	case strings.HasPrefix(mime, "image/"):
		if size > MaxImageSize {
			return fmt.Errorf("%w: max %dMB for images", ErrFileTooLarge, MaxImageSize>>20)
		}
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return fmt.Errorf("%w: only .jpg, .jpeg, .png allowed", ErrUnsupportedFormat)
		}

	case strings.HasPrefix(mime, "video/"):
		if size > MaxVideoSize {
			return fmt.Errorf("%w: max %dMB for videos", ErrFileTooLarge, MaxVideoSize>>20)
		}
		if ext != ".mp4" && ext != ".avi" && ext != ".mov" && ext != ".mkv" {
			return fmt.Errorf("%w: only .mp4, .avi, .mov, .mkv allowed", ErrUnsupportedFormat)
		}

	default:
		return ErrUnsupportedType
	}

	return nil
}
