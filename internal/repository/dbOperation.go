package repository

import (
	"context"
	"fmt"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type workOperation struct {
	db *db
}

func (w *workOperation) Create(ctx context.Context, work domain.Work) (domain.Id, error) {

	var id domain.Id
	var taskId domain.Id
	var duration int
	var resource int

	result, err := w.db.conn.Query(ctx, fmt.Sprintf("INSERT INTO works VALUES($1, $2, $3, $4) RETURNING id, work.Id, work.TaskId, work.Duration, work.Resource"))

	if err != nil {
		fmt.Println("Error occured while inserting a row into database -", err.Error())
		return 0, err
	}

	defer result.Close()

	err = result.Scan(&id, &taskId, &duration, &resource)

	if err != nil {
		fmt.Println("Error occured while inserting a row into database -", err.Error())
		return 0, err
	}

	return id, err
}

func (w *workOperation) Delete(ctx context.Context, id domain.Id) error {
	_, err := w.db.conn.Query(ctx, fmt.Sprintf("DELETE FROM works WHERE id=$1, id"))

	if err != nil {
		fmt.Println("Error occured while deleting a row from database -", err.Error())
		return err
	}

	return err
}

func (w *workOperation) List(ctx context.Context, id domain.Id) (domain.WorkResponse, error) {
	result, err := w.db.conn.Query(ctx, fmt.Sprintf("SELECT id, task_id, duration, password FROM works WHERE id <= $1, id"))

	var response domain.WorkResponse

	if err != nil {
		fmt.Println("Error occured while trying to find a row with id =", id, "-", err.Error())
		return response, err
	}

	result.Next()

	response.Parental = make([]domain.Work, 0)

	var work domain.Work

	for result.Next() {
		err = result.Scan(&work.Id, &work.TaskId, &work.Duration, &work.Resource)

		if err != nil {
			fmt.Println("Error occured while trying to show works rows -", err.Error())
			return response, err
		}

		if id == work.Id {
			response.Main.Id = work.Id
			response.Main.TaskId = work.TaskId
			response.Main.Duration = work.Duration
			response.Main.Resource = work.Resource

		} else {
			response.Parental = append(response.Parental, domain.Work{
				Id:       work.Id,
				TaskId:   work.TaskId,
				Duration: work.Duration,
				Resource: work.Resource,
			})
		}
	}

	return response, err
}
