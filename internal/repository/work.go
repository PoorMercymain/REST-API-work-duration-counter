package repository

import (
	"context"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type work struct {
	db *db
}

func NewWork(db *db) *work {
	return &work{db: db}
}

func (r *work) Create(ctx context.Context, work domain.Work) (domain.Id, error) {

	var id domain.Id

	err := r.db.conn.QueryRow(ctx,
		`INSERT INTO work (task_id, duration, resource, previous_ids) VALUES ($1, $2, $3, $4) RETURNING id`,
		work.TaskId, work.Duration, work.Resource, work.PreviousIds).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *work) Delete(ctx context.Context, id domain.Id) error {
	_, err := r.db.conn.Exec(ctx, `DELETE FROM work WHERE id=$1`, id)

	if err != nil {
		return err
	}

	return err
}

func (r *work) List(ctx context.Context, id domain.Id) (domain.WorkResponse, error) {
	var response domain.WorkResponse
	var tid domain.Id

	err := r.db.conn.QueryRow(ctx, `SELECT task_id FROM work WHERE id = $1`, id).Scan(&tid)

	if err != nil {
		return response, err
	}

	result, err := r.db.conn.Query(ctx,
		`SELECT id, task_id, duration, resource, previous_ids FROM work WHERE id <= $1 and task_id = $2`,
		id, tid)

	if err != nil {
		return response, err
	}

	defer result.Close()

	response.Parental = make([]domain.Work, 0)

	for result.Next() {
		var work domain.Work
		previous := make([]uint32, 0)
		err = result.Scan(&work.Id, &work.TaskId, &work.Duration, &work.Resource, &previous)

		for n := range previous {
			work.PreviousIds[n] = domain.Id(previous[n])
		}

		if err != nil {
			return response, err
		}

		if id == work.Id {
			response.Main = work
		} else {
			response.Parental = append(response.Parental, work)
		}
	}

	return response, err
}

func (r *work) CreateOrUpdateIfNotExists(ctx context.Context, work domain.Work) {
	intPreviousIds := make([]uint32, 0)

	for _, element := range work.PreviousIds {
		intPreviousIds = append(intPreviousIds, uint32(element))
	}

	r.db.conn.QueryRow(ctx,
		`INSERT INTO work (id, task_id, duration, resource, previous_ids) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET task_id = $2, duration = $3, resource = $4, previous_ids = $5`,
		work.Id, work.TaskId, work.Duration, work.Resource, intPreviousIds).Scan()
}
