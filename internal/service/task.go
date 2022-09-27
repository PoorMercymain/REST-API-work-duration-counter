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

func (t task) CreateTask(ctx context.Context, task domain.Task) (domain.Id, error) {
	return t.repo.CreateTask(ctx, task)
}

func (t task) DeleteTask(ctx context.Context, id domain.Id) error {
	return t.repo.DeleteTask(ctx, id)
}

func (t task) GetTask(ctx context.Context, id domain.Id, tid domain.Id) (domain.Task, error) {
	return t.repo.GetTask(ctx, id)
}
