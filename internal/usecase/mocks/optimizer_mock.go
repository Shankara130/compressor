package mocks

type OptimizerMock struct{}

func (m *OptimizerMock) Optimize(input string, output string) error {
	return nil
}
