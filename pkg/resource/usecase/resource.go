package usecase

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
)

type resourceUsecase struct {
	resourceRepo domain.ResourceRepository
}

func NewResourceUsecase(resourceRepo domain.ResourceRepository) resourceUsecase  {
	return resourceUsecase{
		resourceRepo: resourceRepo,
	}
}

func (ru *resourceUsecase) Store(r *domain.Resource) error {
	return ru.resourceRepo.Store(r)
}

func (ru *resourceUsecase) Update(r *domain.Resource) error {
	return ru.resourceRepo.Update(r)
}

func (ru *resourceUsecase) Delete(id uint64) error {
	return ru.resourceRepo.Delete(id)
}

func (ru *resourceUsecase) Find(id uint64) (*domain.Resource, error) {
	return ru.resourceRepo.Find(id)
}

func (ru *resourceUsecase) FindAll(limit, offset int, params map[string]interface{}) ([]*domain.Resource, error) {
	return ru.resourceRepo.FindAll(limit, offset, params)
}