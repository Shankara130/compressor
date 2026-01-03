package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/factory"
)

type OptimizeFileUseCase struct{}

func NewOptimizeFileUseCase() *OptimizeFileUseCase {
	return &OptimizeFileUseCase{}
}

func (u *OptimizeFileUseCase) Execute(file entity.File) error {
	optimizer, err := factory.NewOptimizer(file.MimeType)
	if err != nil {
		return err
	}

	return optimizer.Optimize(file.InputPath, file.OutputPath)
}
