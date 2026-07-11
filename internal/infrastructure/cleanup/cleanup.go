// Package cleanup reaps stale files from the working directories on a
// schedule, so the tmp tree does not grow without bound and files orphaned
// by a crash or restart are eventually removed.
package cleanup

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Sweep removes regular files in each of dirs whose modification time is
// older than maxAge. Subdirectories and missing directories are skipped.
// It is safe to run alongside active uploads/workers: a file being written
// has a fresh modification time and is left alone.
func Sweep(dirs []string, maxAge time.Duration) int {
	cutoff := time.Now().Add(-maxAge)
	removed := 0
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if !info.ModTime().Before(cutoff) {
				continue
			}
			path := filepath.Join(dir, e.Name())
			if err := os.Remove(path); err != nil {
				log.Printf("cleanup: remove %s: %v", path, err)
				continue
			}
			removed++
		}
	}
	return removed
}

// Loop calls Sweep on every interval until ctx is cancelled. A file is
// considered stale once it is older than the interval.
func Loop(ctx context.Context, dirs []string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if n := Sweep(dirs, interval); n > 0 {
				log.Printf("cleanup: removed %d stale file(s)", n)
			}
		}
	}
}
