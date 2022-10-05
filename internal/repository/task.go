package repository

import (
	"context"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

func newScanError(text string) error {
	return &scanError{text}
}

type scanError struct {
	s string
}

func (e *scanError) Error() string {
	return e.s
}

type task struct {
	db *db
}

func NewTask(db *db) *task {
	return &task{db: db}
}

func (r *task) Create(ctx context.Context, task domain.Task) (domain.Id, error) {
	var id domain.Id

	err := r.db.conn.QueryRow(ctx,
		`INSERT INTO task (order_name, start_date) VALUES ($1, $2) RETURNING id`,
		task.OrderName, task.StartDate).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *task) Update(ctx context.Context, id domain.Id, task domain.Task) error {
	_, err := r.db.conn.Exec(ctx, `UPDATE task SET order_name = $1, start_date = $2 WHERE id = $3`, task.OrderName, task.StartDate, id)

	if err != nil {
		return err
	}

	return err
}

func (r *task) Delete(ctx context.Context, id domain.Id) error {
	_, err := r.db.conn.Exec(ctx, `DELETE FROM task WHERE id=$1`, id)

	if err != nil {
		return err
	}

	return err
}

func (r *task) GetTask(ctx context.Context, id domain.Id) (domain.Task, error) {
	var resultTask domain.Task

	result, err := r.db.conn.Query(ctx,
		`SELECT id, order_name, start_date FROM task WHERE id = $1`,
		id)

	if err != nil {
		return resultTask, err
	}

	defer result.Close()

	if result.Next() {
		err = result.Scan(&resultTask.Id, &resultTask.OrderName, &resultTask.StartDate)
		if err != nil {
			return resultTask, err
		}
		return resultTask, err
	} else {
		scanErr := newScanError("Error occured while scaning sql row")
		return resultTask, scanErr
	}
}
