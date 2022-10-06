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

type Node interface {
	AddNext() Node
	MoveTo() []Node
}

type TreeNode struct {
	children []Node
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
		scanErr := newScanError("Error occurred while scanning sql row")
		return resultTask, scanErr
	}
}

func (r *task) ListWorksOfTask(ctx context.Context, id domain.Id) ([]domain.Work, error) {
	var resultWorks = make([]domain.Work, 0)

	result, err := r.db.conn.Query(ctx,
		`SELECT id, task_id, duration, resource FROM work WHERE task_id = $1`,
		id)

	if err != nil {
		return resultWorks, err
	}

	defer result.Close()

	var currentWork domain.Work

	for result.Next() {
		err = result.Scan(&currentWork.Id, &currentWork.TaskId, &currentWork.Duration, &currentWork.Resource)
		if err != nil {
			return resultWorks, err
		}

		resultWorks = append(resultWorks, currentWork)
	}
	return resultWorks, err
}

func (r *task) FindRoot(ctx context.Context, works []domain.Work) domain.Work {
	var rootId domain.Id

	for _, currentNode := range works {
		if rootId == currentNode.Id {
			return currentNode
		}
	}
	return domain.Work{Id: 0, TaskId: 0, Duration: 0, Resource: 0}
}

func (r *task) FindLeafs(ctx context.Context, works domain.Work) []domain.Work {
	return make([]domain.Work, 0) //placeholder
}
