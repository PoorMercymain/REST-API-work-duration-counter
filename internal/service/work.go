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

func (w work) CreateWork(ctx context.Context, work domain.Work) (domain.Id, error) {
	return w.repo.CreateWork(ctx, work)
}

func (w work) DeleteWork(ctx context.Context, id domain.Id) error {
	return w.repo.DeleteWork(ctx, id)
}

func (w work) ListWork(ctx context.Context, id domain.Id, tid domain.Id) (domain.WorkResponse, error) {
	return w.repo.ListWork(ctx, id, tid)
}

func (w work) Count(workResponse domain.WorkResponse) int {
	result := workResponse.Main.Duration
	for _, parentalElement := range workResponse.Parental {
		result += parentalElement.Duration
	}

	return result
}
