package optimizer

import "os/exec"

type VideoOptimizer struct{}

func NewVideoOptimizer() *VideoOptimizer {
	return &VideoOptimizer{}
}

func (o *VideoOptimizer) Optimize(input, output string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", input,
		"-vcodec", "libx264",
		"-crf", "28",
		"-preset", "slow",
		output,
	)
	return cmd.Run()
}
