package optimizer

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Shankara130/compressor/internal/domain/service"
)

type VideoOptimizer struct{}

func NewVideoOptimizer() *VideoOptimizer {
	return &VideoOptimizer{}
}

func (o *VideoOptimizer) Optimize(input, output string, progress service.ProgressFunc) error {
	_ = os.Remove(output)

	// Total duration is needed to turn ffmpeg's position into a percentage.
	// If it can't be determined, we simply don't report progress.
	duration := probeDuration(input)

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", input,
		"-movflags", "+faststart",
		"-vcodec", "libx264",
		"-crf", "28",
		"-preset", "slow",
		"-acodec", "aac",
		"-progress", "pipe:1", // structured progress key/values to stdout
		"-nostats", // keep stderr to errors only
		output,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg start: %w", err)
	}

	// Parse ffmpeg's progress lines and report position within the file.
	progressDone := make(chan struct{})
	go func() {
		defer close(progressDone)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "out_time=") {
				continue
			}
			if duration <= 0 || progress == nil {
				continue
			}
			secs, ok := parseSeconds(strings.TrimPrefix(line, "out_time="))
			if !ok {
				continue
			}
			pct := min(99, max(0, int(secs/duration*100))) // never claim completion
			progress(pct)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("ffmpeg progress: %v", err)
		}
	}()

	waitErr := cmd.Wait()
	<-progressDone // make sure no progress callback fires after completion

	if waitErr != nil {
		return fmt.Errorf("ffmpeg failed: %s", stderr.String())
	}
	return nil
}

// probeDuration returns the media duration in seconds via ffprobe, or 0 if it
// cannot be determined.
func probeDuration(input string) float64 {
	out, err := exec.Command("ffprobe", "-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", input).Output()
	if err != nil {
		return 0
	}
	d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil || d <= 0 {
		return 0
	}
	return d
}

// parseSeconds parses ffmpeg's time format (HH:MM:SS or HH:MM:SS.microseconds).
func parseSeconds(s string) (float64, bool) {
	main, frac, hasFrac := strings.Cut(s, ".")
	parts := strings.Split(main, ":")
	if len(parts) != 3 {
		return 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	sec, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, false
	}
	total := float64(h*3600 + m*60 + sec)
	if hasFrac {
		if us, err := strconv.Atoi(frac); err == nil {
			total += float64(us) / 1e6
		}
	}
	return total, true
}
