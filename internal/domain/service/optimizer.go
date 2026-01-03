package service

type Optimizer interface {
	Optimize(inputPath string, outputPath string) error
}
