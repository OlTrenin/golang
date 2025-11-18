package service

import (
	"fmt"

	"number-service/internal/domain"
)

type numberService struct {
	repo domain.NumberRepository
}

func NewNumberService(repo domain.NumberRepository) domain.NumberService {
	return &numberService{repo: repo}
}

func (s *numberService) AddNumber(value int) ([]int, error) {
	if err := s.repo.Save(value); err != nil {
		return nil, fmt.Errorf("failed to save number: %w", err)
	}

	numbers, err := s.repo.GetAllSorted()
	if err != nil {
		return nil, fmt.Errorf("failed to get sorted numbers: %w", err)
	}

	return numbers, nil
}

var _ domain.NumberService = (*numberService)(nil)
