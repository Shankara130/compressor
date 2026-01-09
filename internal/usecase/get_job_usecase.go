package usecase

import (
	"github.com/Shankara130/compressor/internal/domain/entity"
	"github.com/Shankara130/compressor/internal/domain/service"
)

type GetJobUseCase struct {
	Repository service.JobRepository
}

func NewGetJobUseCase(r service.JobRepository) *GetJobUseCase {
	return &GetJobUseCase{Repository: r}
}

func (u *GetJobUseCase) Execute(id string) (entity.Job, error) {
	return u.Repository.GetByID(id)
}
