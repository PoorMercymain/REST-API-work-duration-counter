package domain

import (
	"context"
	"time"
)

type Task struct {
	Id        Id        `json:"id"`
	OrderName string    `json:"order_name"`
	StartDate time.Time `json:"start_date"`
}

type TaskRepository interface {
	Create(ctx context.Context, task Task) (Id, error)
	Update(ctx context.Context, id Id, task Task) error
	Delete(ctx context.Context, id Id) error
	GetTask(ctx context.Context, id Id) (Task, error)
}

// TaskService тут бизнес-логика
type TaskService interface {
	Create(ctx context.Context, task Task) (Id, error)
	Update(ctx context.Context, id Id, task Task) error
	Delete(ctx context.Context, id Id) error
}
