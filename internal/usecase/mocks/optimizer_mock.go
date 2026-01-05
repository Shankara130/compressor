package mocks

type OptimizerMock struct {
	Called bool
	Err    error
}

func (m *OptimizerMock) Optimize(input string, output string) error {
	m.Called = true
	return m.Err
}
