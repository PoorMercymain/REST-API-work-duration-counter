package service

import (
	"context"
	"fmt"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type work struct {
	repo domain.WorkRepository
}

func NewWork(repo domain.WorkRepository) *work {
	return &work{repo: repo}
}

func (s *work) Create(ctx context.Context, work domain.Work) (domain.Id, error) {
	return s.repo.Create(ctx, work)
}

func (s *work) Delete(ctx context.Context, id domain.Id) error {
	return s.repo.Delete(ctx, id)
}

func (s *work) List(ctx context.Context, id domain.Id) (domain.WorkResponse, error) {
	return s.repo.List(ctx, id)
}

func (s *work) CreateTestWorks(ctx context.Context) {
	for i := 0; i < 100000; i++ {
		firstAndThird := make([]domain.Id, 0)
		firstAndThird = append(firstAndThird, domain.Id(1+(i*26)))
		firstAndThird = append(firstAndThird, domain.Id(3+(i*26)))
		first := make([]domain.Id, 0)
		first = append(first, domain.Id(1+(i*26)))
		third := make([]domain.Id, 0)
		third = append(third, domain.Id(3+(i*26)))
		secondThirdAndFourth := make([]domain.Id, 0)
		secondThirdAndFourth = append(secondThirdAndFourth, domain.Id(2+(i*26)))
		secondThirdAndFourth = append(secondThirdAndFourth, domain.Id(3+(i*26)))
		secondThirdAndFourth = append(secondThirdAndFourth, domain.Id(4+(i*26)))

		fmt.Println(i+1, "iter")
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(1 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 45000, Resource: 7})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(2 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 2500, Resource: 6})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(3 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 500, Resource: 8})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(4 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 40000, Resource: 5, PreviousIds: third})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(5 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 45000, Resource: 7, PreviousIds: firstAndThird})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(6 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 4000, Resource: 8, PreviousIds: first})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(7 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 3000, Resource: 2, PreviousIds: secondThirdAndFourth})
		s.repo.CreateOrUpdateIfNotExists(ctx, domain.Work{Id: domain.Id(8 + (i * 26)), TaskId: domain.Id(i + 1), Duration: 1000, Resource: 9})
	}
}

// TODO: private
func (w work) Count(workResponse domain.WorkResponse) int {
	result := workResponse.Main.Duration
	for _, parentalElement := range workResponse.Parental {
		result += parentalElement.Duration
	}

	return result
}
