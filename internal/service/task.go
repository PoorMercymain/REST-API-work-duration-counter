package service

import (
	"context"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type task struct {
	repo domain.TaskRepository
}

func NewTask(repo domain.TaskRepository) *task {
	return &task{repo: repo}
}

func (s *task) Create(ctx context.Context, task domain.Task) (domain.Id, error) {
	return s.repo.Create(ctx, task)
}

func (s *task) Update(ctx context.Context, id domain.Id, task domain.Task) error {
	return s.repo.Update(ctx, id, task)
}

func (s *task) Delete(ctx context.Context, id domain.Id) error {
	return s.repo.Delete(ctx, id)
}

func (s *task) GetTask(ctx context.Context, id domain.Id) (domain.Task, error) {
	return s.repo.GetTask(ctx, id)
}

func (s *task) ListWorksOfTask(ctx context.Context, id domain.Id) ([]domain.Work, error) {
	return s.repo.ListWorksOfTask(ctx, id)
}
