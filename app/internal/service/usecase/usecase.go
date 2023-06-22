package usecase

import (
	serviceRep "TP2023_DB/app/internal/service/repository"
	"TP2023_DB/app/models"
)

type UseCaseI interface {
	ClearData() error
	SelectStatus() (*models.ServiceStatus, error)
}

type useCase struct {
	serviceRepository serviceRep.RepositoryI
}

func New(serviceRepository serviceRep.RepositoryI) UseCaseI {
	return &useCase{
		serviceRepository: serviceRepository,
	}
}

func (uc *useCase) ClearData() error {
	err := uc.serviceRepository.ClearData()
	return err
}

func (uc *useCase) SelectStatus() (*models.ServiceStatus, error) {
	status, err := uc.serviceRepository.SelectStatus()
	if err != nil {
		return nil, err
	}

	return status, nil
}
