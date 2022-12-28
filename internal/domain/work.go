package domain

import "context"

type Work struct {
	Id          Id   `json:"id"`
	TaskId      Id   `json:"task_id"`
	Duration    int  `json:"duration"`
	Resource    int  `json:"resource"`
	PreviousIds []Id `json:"previous_ids"`
}

type WorkResponse struct {
	Main     Work   `json:"main"`
	Parental []Work `json:"parental"`
}

type WorkRepository interface {
	Create(ctx context.Context, work Work) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id) (WorkResponse, error)
	CreateOrUpdateIfNotExists(ctx context.Context, work Work)
}

type WorkService interface {
	Count(workResponse WorkResponse) int
	Create(ctx context.Context, task Work) (Id, error)
	Delete(ctx context.Context, id Id) error
	List(ctx context.Context, id Id) (WorkResponse, error)
	CreateTestWorks(ctx context.Context)
}
