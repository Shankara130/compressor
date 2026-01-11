package optimizer

import (
	"fmt"
	"os"
	"os/exec"
)

type VideoOptimizer struct{}

func NewVideoOptimizer() *VideoOptimizer {
	return &VideoOptimizer{}
}

func (o *VideoOptimizer) Optimize(input, output string) error {
	_ = os.Remove(output)

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", input,
		"-movflags", "+faststart",
		"-vcodec", "libx264",
		"-crf", "28",
		"-preset", "slow",
		"-acodec", "aac",
		output,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %s", string(out))
	}

	return nil
}
