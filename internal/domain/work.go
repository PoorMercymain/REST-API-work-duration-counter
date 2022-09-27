package domain

import "context"

type Work struct {
	Id       Id  `json:"id"`
	TaskId   Id  `json:"task_id"`
	Duration int `json:"duration"`
	Resource int `json:"resource"`
}

type WorkResponse struct {
	Main     Work   `json:"main"`
	Parental []Work `json:"parental"`
}

type WorkRepository interface {
	Create(ctx context.Context, work Work) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id, tid Id) (WorkResponse, error)
}

type WorkService interface {
	//здесь тоже бизнес-логика
	Count(workResponse WorkResponse) int

	Create(ctx context.Context, task Work) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id, tid Id) (WorkResponse, error)
}
