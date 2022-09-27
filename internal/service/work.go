package service

import (
	"context"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type work struct {
	repo domain.WorkRepository
}

func NewWork(repo domain.WorkRepository) *work {
	return &work{repo: repo}
}

func (s *work) Create(ctx context.Context, work domain.Work) (domain.Id, error) {
	return s.repo.Create(ctx, work)
}

func (s *work) Delete(ctx context.Context, id domain.Id) error {
	return s.repo.Delete(ctx, id)
}

func (s *work) List(ctx context.Context, id domain.Id, tid domain.Id) (domain.WorkResponse, error) {
	return s.repo.List(ctx, id, tid)
}

// TODO: private
func (w work) Count(workResponse domain.WorkResponse) int {
	result := workResponse.Main.Duration
	for _, parentalElement := range workResponse.Parental {
		result += parentalElement.Duration
	}

	return result
}
