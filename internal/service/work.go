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

func (w work) Create(ctx context.Context, work domain.Work) (domain.Id, error) {
	return w.repo.Create(ctx, work)
}

func (w work) Delete(ctx context.Context, id domain.Id) error {
	//TODO implement me
	panic("implement me")
}

func (w work) List(ctx context.Context, id domain.Id, tid domain.Id) (domain.WorkResponse, error) {
	//TODO implement me
	panic("implement me")
}
