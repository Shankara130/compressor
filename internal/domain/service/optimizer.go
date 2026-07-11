package service

// ProgressFunc reports an optimizer's own progress as a value in [0, 100].
// It may be called from any goroutine; implementations must ignore a nil
// callback.
type ProgressFunc func(percent int)

type Optimizer interface {
	Optimize(inputPath string, outputPath string, progress ProgressFunc) error
}
