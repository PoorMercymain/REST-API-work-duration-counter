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
		`SELECT id, task_id, duration, resource, previous_ids FROM work WHERE task_id = $1`,
		id)

	if err != nil {
		return resultWorks, err
	}

	defer result.Close()

	var currentWork domain.Work

	var previousIdsBuffer = make([]uint32, 0)
	for result.Next() {

		err = result.Scan(&currentWork.Id, &currentWork.TaskId, &currentWork.Duration, &currentWork.Resource, &previousIdsBuffer)
		if err != nil {
			return resultWorks, err
		}

		currentWork.PreviousIds = make([]domain.Id, 0)

		for _, previousId := range previousIdsBuffer {
			currentWork.PreviousIds = append(currentWork.PreviousIds, domain.Id(previousId))
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

func (r *task) dropElementFromWorksSlice(index int, someSlice []domain.Work) []domain.Work {
	return append(someSlice[:index], someSlice[index+1:]...)
}

func (r *task) UpdateOrCreateIfNotExists(ctx context.Context, task domain.Task) error {
	var id domain.Id

	err := r.db.conn.QueryRow(ctx,
		`INSERT INTO task (id, order_name, start_date) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET order_name = $2, start_date = $3`,
		task.Id, task.OrderName, task.StartDate).Scan(&id)

	if err != nil {
		return err
	}

	return err
} //INSERT INTO table (a,b,c) VALUES (4,5,6) ON DUPLICATE KEY UPDATE c=9;

/*func (r *task) FindLeafs(ctx context.Context, works []domain.Work) []domain.Work {
	//находим root
	root := r.FindRoot(ctx, works)
	//начальная инициализация слайса последних узлов
	lastNodes := make([]domain.Work, 0)

	//цикл по всем работам
	for counter, element := range works {
		//если id элемента равно id root-а (т.е. если текущий элемент - root)
		if element.Id == root.Id {
			//удаляем элемент из слайса работ
			works = r.dropElementFromWorksSlice(counter, works)
			//добавляем в слайс последних узлов root
			lastNodes = append(lastNodes, root)
			//выходим из цикла
			break
		}
	}

	//цикл по всем work-ам
	for _, element := range works {
		//начальная инициализация слайса проверенных работ (т.е. тех, у которых есть предшествующие работы)
		checkedNodes := make([]domain.Work, 0)
		//инициализация слайса работ, которые еще нужно проверить (т.е. те, которые являются дочерними по отношению к текущим)
		nodesToCheck := make([]domain.Work, 0)
		//цикл по всем последним узлам
		for _, node := range lastNodes {
			//если среди последних узлов есть такой, чей id == parentId текущей работы, добавляем текущую работу в слайс работ, требующих проверки, а последний узел - в проверенные работы
			if element.PreviousIds == node.Id {
				checkedNodes = append(checkedNodes, node)
				nodesToCheck = append(nodesToCheck, element)
			}
		}

		//цикл по всем проверенным узлам
		for _, currentElement := range checkedNodes {
			//цикл по всем последним узлам
			for i, currentLastNodesElement := range lastNodes {
				//если проверенный узел есть среди последних, то он не последний, так что удаляем его
				if currentElement == currentLastNodesElement {
					lastNodes = r.dropElementFromWorksSlice(i, lastNodes)
				}
			}
		}
		//из слайса проверенных работ убираем все элементы
		checkedNodes = make([]domain.Work, 0)
		//цикл по всем работам
		for i, worksElement := range works {
			//цикл по всем элементам, требующим проверки
			for _, toCheckElement := range nodesToCheck {
				//если среди всех работ есть проверяемая работа, удаляем ее из слайса работ и добавляем в слайс последних работ
				if worksElement.Id == toCheckElement.Id {
					works = r.dropElementFromWorksSlice(i, works)
					lastNodes = append(lastNodes, toCheckElement)
				}
			}
		}
	}
	//возвращаем слайс последних работ
	return lastNodes
}*/
