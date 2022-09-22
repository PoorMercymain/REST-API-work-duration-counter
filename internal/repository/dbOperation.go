package repository

import (
	"context"
	"fmt"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
	"github.com/jackc/pgx/v4"
)

type workOperation struct {
	conn *pgx.Conn
}

func (w *workOperation) Create(ctx context.Context, work domain.Work) (domain.Id, error) {
	result, err := w.conn.Query(ctx, fmt.Sprintf("INSERT INTO works VALUES(%d, %d, %d, %d)", work.Id, work.TaskId, work.Duration, work.Resource))

	if err != nil {
		fmt.Println("Error occured while inserting a row into database -", err.Error())
		return 0, err
	}

	var id domain.Id
	var taskId domain.Id
	var duration int
	var resource int

	result.Next()
	err = result.Scan(&id, &taskId, &duration, &resource)

	if err != nil {
		fmt.Println("Error occured while trying to get inserted row id -", err.Error())
		return 0, err
	}

	return id, err
}

func (w *workOperation) Delete(ctx context.Context, id domain.Id) error {
	_, err := w.conn.Query(ctx, fmt.Sprintf("DELETE FROM works WHERE id=%d", id))

	if err != nil {
		fmt.Println("Error occured while deleting a row from database -", err.Error())
		return err
	}

	return err
}

func (w *workOperation) List(ctx context.Context, id domain.Id) (domain.WorkResponse, error) {
	result, err := w.conn.Query(ctx, fmt.Sprintf("SELECT id, task_id, duration, password FROM works WHERE id <= %d", id))

	var response domain.WorkResponse

	if err != nil {
		fmt.Println("Error occured while trying to find a row with id =", id, "-", err.Error())
		return response, err
	}

	result.Next()

	response.Parental = make([]domain.Work, 0)

	var workId domain.Id
	var taskId domain.Id
	var duration int
	var resource int

	for {
		err = result.Scan(&workId, &taskId, &duration, &resource)

		if err != nil {
			fmt.Println("Error occured while trying to show works rows -", err.Error())
			return response, err
		}

		if id == workId {
			response.Main = domain.Work{
				Id:       workId,
				TaskId:   taskId,
				Duration: duration,
				Resource: resource,
			}
		} else {
			response.Parental = append(response.Parental, domain.Work{
				Id:       workId,
				TaskId:   taskId,
				Duration: duration,
				Resource: resource,
			})
		}

		if !(result.Next()) {
			break
		}
	}

	return response, err
}
