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
	CreateTask(ctx context.Context, task Task) (Id, error)
	UpdateTask(ctx context.Context, id Id, task Task) error
	DeleteTask(ctx context.Context, id Id) error
	GetTask(ctx context.Context, id Id) (Task, error)
}

// TaskService тут бизнес-логика
type TaskService interface {
	CreateTask(ctx context.Context, task Task) (Id, error)
	UpdateTask(ctx context.Context, id Id, task Task) error
	DeleteTask(ctx context.Context, id Id) error
}
