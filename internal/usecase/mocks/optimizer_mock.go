package mocks

import "github.com/Shankara130/compressor/internal/domain/service"

type OptimizerMock struct{}

func (m *OptimizerMock) Optimize(input string, output string, _ service.ProgressFunc) error {
	return nil
}
