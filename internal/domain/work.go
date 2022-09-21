package domain

import "context"

type Work struct {
	Id       Id  `json:"id"`
	TaskId   Id  `json:"task_id"`
	Duration int `json:"duration"`
	Resource int `json:"resouce"`
}

type WorkResponse struct {
	Main     Work   `json:"main"`
	Parental []Work `json:"parental"`
}

type WorkRepository interface {
	Create(ctx context.Context, task Task) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id) (WorkResponse, error)
}

type WorkService interface {
	//здесь тоже бизнес-логика
	Create(ctx context.Context, task Task) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id) (WorkResponse, error)
}
