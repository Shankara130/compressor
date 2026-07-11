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

// MaxSizeForMime returns the maximum allowed upload size for the given content
// type, or 0 if the type is not supported. Single source of truth for size
// limits, used both by ValidateFile and by the streaming upload handler.
func MaxSizeForMime(mime string) int64 {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return MaxImageSize
	case strings.HasPrefix(mime, "video/"):
		return MaxVideoSize
	case mime == "application/pdf":
		return MaxPDFSize
	default:
		return 0
	}
}

func ValidateFile(mime string, size int64, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	maxSize := MaxSizeForMime(mime)
	if maxSize == 0 {
		return ErrUnsupportedType
	}
	if size > maxSize {
		return fmt.Errorf("%w: max %dMB", ErrFileTooLarge, maxSize>>20)
	}

	switch {
	case strings.HasPrefix(mime, "image/"):
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return fmt.Errorf("%w: only .jpg, .jpeg, .png allowed", ErrUnsupportedFormat)
		}

	case strings.HasPrefix(mime, "video/"):
		if ext != ".mp4" && ext != ".avi" && ext != ".mov" && ext != ".mkv" {
			return fmt.Errorf("%w: only .mp4, .avi, .mov, .mkv allowed", ErrUnsupportedFormat)
		}

	case mime == "application/pdf":
		if ext != ".pdf" {
			return fmt.Errorf("%w: only .pdf allowed", ErrUnsupportedFormat)
		}
	}

	return nil
}
