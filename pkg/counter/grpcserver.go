package counter

import (
	"context"
	"fmt"
	api "github.com/PoorMercymain/REST-API-work-duration-counter/pkg/api/api/proto"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type GRPCServer struct{}

func (G GRPCServer) Count(ctx context.Context, request *api.CountRequest) (*api.CountResponse, error) {
	duration, path, err := CountAllDuration()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strconv.Itoa(duration), path)
	fmt.Println(request.CountNeeded, "<- sent from server")
	return &api.CountResponse{Result: "110"}, nil
}

func (G GRPCServer) MustEmbedUnimplementedCounterServer() {

	panic("this method is redundant")
}

func RedisConnect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
func RedisSet(rdb *redis.Client, key string, value string) {
	var ctx = context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
func RedisGet(rdb *redis.Client, key string) string {
	var ctx = context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

/*func (s *task) CreateTestTasks(ctx context.Context) error {
	//var workRepo WorkRepository
	var err error

	//var initTask Task
	generationChan := make([]chan int, 0)
	for j := 0; j < 10; j++ {
		var currentChan chan int
		generationChan = append(generationChan, currentChan)
		go func() {
			for i := 0; i < 100000/4; i++ {
				err = s.repo.UpdateOrCreateIfNotExists(ctx, Task{Id: Id((i * j) + 1), OrderName: fmt.Sprintf("Task number %d", (i*j)+1), StartDate: time.Now()})
			}
			currentChan <- 1
		}()
	}

	for _, channel := range generationChan {
		<-channel
	}

	return err
}

func (s *task) CountDuration(ctx context.Context, id Id) (uint64, error) {
	allWorksOfTaskSlice, err := s.repo.ListWorksOfTask(ctx, id)

	rand.Seed(time.Now().Unix())

	if err != nil {
		return 0, err
	}

	doneWorksSlice := make([]Work, 0)

	undoneWorksSlice := make([]Work, 0)

	availableWorksSlice := make([]Work, 0)

	inProgressWorksSlice := make([]Work, 0)

	resources := 10

	toDropAvailableIds := make([]Id, 0)

	toDropInProgressIds := make([]Id, 0)

	choosen := make([]Work, 0)

	var minDuration int

	var result int

	for _, workOfTask := range allWorksOfTaskSlice {
		undoneWorksSlice = append(undoneWorksSlice, workOfTask)
	}
	var u int
	for {
		u += 1
		fmt.Println(u)
		for _, workOfTask := range allWorksOfTaskSlice {
			//если ворка нет среди выполненных, доступных для выполнения или находящихся в обработке - проверяем дальше
			if (findWorkById(doneWorksSlice, workOfTask.Id) == -1) && (findWorkById(availableWorksSlice, workOfTask.Id) == -1) && (findWorkById(inProgressWorksSlice, workOfTask.Id) == -1) {
				var doneForTheWork int
				for _, prevId := range workOfTask.PreviousIds {
					//если работа из списка предшествующих есть среди выполненных, увеличиваем счетчик
					if findWorkById(doneWorksSlice, prevId) != -1 {
						doneForTheWork += 1
					}
				}
				if len(workOfTask.PreviousIds) == doneForTheWork { //если значение счетчика равно количеству предшествующих работ, значит работу можно выполнять - добавляем ее в available
					availableWorksSlice = append(availableWorksSlice, workOfTask)
				}
			}
		}

		fmt.Println(availableWorksSlice, "available/")

		//сортируем слайс available по возрастанию
		sort.Slice(availableWorksSlice, func(i, j int) bool {
			return availableWorksSlice[i].Duration < availableWorksSlice[j].Duration
		})

		fmt.Println(availableWorksSlice, "/available")

		var ran int
		for _ = range availableWorksSlice {
			randomAvailable := availableWorksSlice[rand.Intn(len(availableWorksSlice))]
			for findWorkById(choosen, randomAvailable.Id) != -1 {
				ran = rand.Intn(len(availableWorksSlice))
				randomAvailable = availableWorksSlice[ran]
				fmt.Println(ran, "randomness")
			}

			//если позволяют ресурсы - добавляем работу в слайс "в обработке", уменьшаем доступные ресурсы на ресурсы работы и добавляем id в слайс для последующего удаления
			if resources >= randomAvailable.Resource {
				resources -= randomAvailable.Resource
				toDropAvailableIds = append(toDropAvailableIds, randomAvailable.Id)
				inProgressWorksSlice = append(inProgressWorksSlice, randomAvailable)
				choosen = append(choosen, randomAvailable)
				fmt.Println(ran, "random")
			}
		}
		//удаляем элементы, добавленные в обработку из доступных
		for _, element := range toDropAvailableIds {
			availableWorksSlice = dropElementFromWorkSliceById(availableWorksSlice, element)
		}
		//очищаем слайс id
		toDropAvailableIds = make([]Id, 0)

		//находим минимальную длительность среди выполняемых работ, а также индекс соответствующего элемента
		_, minDuration = findMinDurationAndIndex(inProgressWorksSlice)
		fmt.Println(inProgressWorksSlice, "pam")

		//добавляем найденную продолжительность в результат
		result += minDuration

		fmt.Println(minDuration, "duration")

		fmt.Println(inProgressWorksSlice, "in progress")

		//убавляем длительность всех работ на найденную минимальную
		for i := range inProgressWorksSlice {
			inProgressWorksSlice[i].Duration -= minDuration
			//если продолжительность == 0 - работа выполнена и ее надо удалить
			if inProgressWorksSlice[i].Duration == 0 {
				toDropInProgressIds = append(toDropInProgressIds, inProgressWorksSlice[i].Id)
				doneWorksSlice = append(doneWorksSlice, inProgressWorksSlice[i])
				resources += inProgressWorksSlice[i].Resource
			}
		}

		fmt.Println(toDropInProgressIds, "->")

		//удаляем выполненные работы из слайсов "в обработке" и undone
		for _, currentId := range toDropInProgressIds {
			inProgressWorksSlice = dropElementFromWorkSliceById(inProgressWorksSlice, currentId)

			undoneWorksSlice = dropElementFromWorkSliceById(undoneWorksSlice, currentId)
		}

		fmt.Println(undoneWorksSlice, "<-")

		//очищаем слайс удаляемых из обработки работ
		toDropInProgressIds = make([]Id, 0)

		//когда нету невыполненных работ - возвращаем вычисленное значение
		if len(undoneWorksSlice) == 0 {
			return uint64(result), nil
		}
	}

	/*allWorksSlice, err := s.repo.ListWorksOfTask(ctx, id)
	worksSlice := allWorksSlice
	allWorksNotDone := allWorksSlice

	fmt.Println("all works initial:", allWorksSlice, "works initial:", worksSlice)

	if err != nil {
		return 0, err
	}

	worksAvailable := make([]Work, 0)

	for _, currentWork := range allWorksSlice {
		if len(currentWork.PreviousIds) == 0 {
			worksAvailable = append(worksAvailable, currentWork)
		}
	}

	fmt.Println(worksAvailable, "debug1") //debug

	sort.Slice(worksAvailable, func(i, j int) bool {
		return worksAvailable[i].Duration < worksAvailable[j].Duration
	})

	fmt.Println(worksAvailable, "debug") //debug

	resources := 10
	worksInProgress := make([]Work, 0)
	worksDone := make([]Work, 0)

	for _, currentWork := range worksAvailable {
		if resources == 0 {
			fmt.Println("not enough resources")
			break
		}

		if !(currentWork.Resource > resources) {
			fmt.Println("before adding:", worksInProgress)
			worksInProgress = append(worksInProgress, currentWork)
			resources -= currentWork.Resource
			fmt.Println("after:", worksInProgress, "resources:", resources)
		}
	}

	var minDurationIndex int
	var minDuration int

	worksSliceInitialLength := len(worksSlice)
	fmt.Println(worksSliceInitialLength, "init len of works")

	fmt.Println(worksSliceInitialLength, "debug2")

	elementsIdsToDrop := make([]Id, 0)

	var result int

	for i := 0; i < worksSliceInitialLength; i++ {
		fmt.Println("now counting")
		minDurationIndex, minDuration = findMinDurationAndIndex(worksInProgress)
		fmt.Println("min duration:", minDuration, "index:", minDurationIndex)
		fmt.Println("works in progress:", worksInProgress)
		for _, workInProgress := range worksInProgress {
			fmt.Println("current work:", workInProgress)
			fmt.Println("min duration work:", worksInProgress[minDurationIndex])
			if worksInProgress[minDurationIndex].Duration == workInProgress.Duration || workInProgress.Duration == 0 {
				fmt.Println("before adding new elements to drop:", elementsIdsToDrop)
				elementsIdsToDrop = append(elementsIdsToDrop, workInProgress.Id)
				worksDone = append(worksDone, workInProgress)
				allWorksNotDone = dropElementFromWorkSliceById(allWorksNotDone, workInProgress.Id)
				worksAvailable = dropElementFromWorkSliceById(worksAvailable, workInProgress.Id)
				fmt.Println("not done:", allWorksNotDone)
				fmt.Println(worksDone, "done")
				resources += workInProgress.Resource
			}
		}

		workIndex := -1
		for _, workToDropId := range elementsIdsToDrop {
			fmt.Println("works:", worksInProgress, "to drop:", workToDropId)
			workIndex = findWorkById(worksInProgress, workToDropId)
			if workIndex != -1 {
				fmt.Println("dropping", workIndex)
				if workIndex != len(worksInProgress)-1 {
					worksInProgress = append(worksInProgress[:workIndex], worksInProgress[workIndex+1:]...)
				} else {
					worksInProgress = worksInProgress[:workIndex]
				}
				fmt.Println("after drop:", worksInProgress)
			}
		}

		elementsIdsToDrop = make([]Id, 0)

		for ind := range worksInProgress {
			worksInProgress[ind].Duration -= minDuration
		}

		fmt.Println("works in progress after duration:", worksInProgress)
		fmt.Println("result was:", result)
		result += minDuration
		fmt.Println("result is:", result)

		fmt.Println("all works:", allWorksSlice)
		for _, currentWork := range allWorksNotDone {
			count := 0
			for j := range currentWork.PreviousIds {
				if findWorkById(worksDone, currentWork.PreviousIds[j]) != -1 { //если все previous id есть в done - добавляем
					count += 1
				}
			}
			fmt.Println("previous ids done found for", currentWork, ":", count)
			if count == len(currentWork.PreviousIds) && (findWorkById(worksAvailable, currentWork.Id) == -1) { //если нашлось столько же элементов, сколько и в previous id
				worksAvailable = append(worksAvailable, currentWork)
				fmt.Println(currentWork, "added to available")
				fmt.Println("available:", worksAvailable)
			}

		}

		sort.Slice(worksAvailable, func(i, j int) bool {
			return worksAvailable[i].Duration < worksAvailable[j].Duration
		})

		fmt.Println("after sorting available:", worksAvailable)

		for _, currentWork := range worksAvailable {
			if resources == 0 {
				break
			}

			if !(currentWork.Resource > resources) {
				worksInProgress = append(worksInProgress, currentWork)
				resources -= currentWork.Resource
			}
		}

		fmt.Println(len(worksDone), "works done len")
		if len(allWorksNotDone) == 0 {
			return uint64(result), nil
		}
	}
	return 403, nil*/
//}

type Id uint32

type Task struct {
	Id        Id        `json:"id"`
	OrderName string    `json:"order_name"`
	StartDate time.Time `json:"start_date"`
}

type Work struct {
	Id          Id   `json:"id"`
	TaskId      Id   `json:"task_id"`
	Duration    int  `json:"duration"`
	Resource    int  `json:"resource"`
	PreviousIds []Id `json:"previous_ids"`
}

func findMinValueAndIndex(slice []int) (int, int) {
	min := 0
	index := 0

	fmt.Println(slice, "slice")

	for i, value := range slice {
		if value < min || i == 0 {
			min = value
			index = i
		}
	}

	fmt.Println(min, "min")
	return min, index
}

func CountAllDuration() (int, string, error) {
	type resultStruct struct {
		Duration  int
		WorksPath []Work
	}

	amount := 1000
	allTasks := make([]Task, 0)
	for v := 0; v < amount; v++ {
		allTasks = append(allTasks, Task{Id: Id(v + 1), OrderName: fmt.Sprintf("Task number %d", v+1), StartDate: time.Now()})
	}

	worksBuffer := make([]Work, 0)
	goDoneChans := make([]chan bool, 0)
	result := make([]resultStruct, 0)
	for v := 0; v < amount; v++ {
		firstAndThird := make([]Id, 0)
		firstAndThird = append(firstAndThird, Id(1+(v*26)))
		firstAndThird = append(firstAndThird, Id(3+(v*26)))
		first := make([]Id, 0)
		first = append(first, Id(1+(v*26)))
		third := make([]Id, 0)
		third = append(third, Id(3+(v*26)))
		secondThirdAndFourth := make([]Id, 0)
		secondThirdAndFourth = append(secondThirdAndFourth, Id(2+(v*26)))
		secondThirdAndFourth = append(secondThirdAndFourth, Id(3+(v*26)))
		secondThirdAndFourth = append(secondThirdAndFourth, Id(4+(v*26)))

		//fmt.Println(v+1, "iter")
		worksBuffer = append(worksBuffer, Work{Id: Id(1 + (v * 26)), TaskId: Id(v + 1), Duration: 45000, Resource: 7})
		worksBuffer = append(worksBuffer, Work{Id: Id(2 + (v * 26)), TaskId: Id(v + 1), Duration: 2500, Resource: 6})
		worksBuffer = append(worksBuffer, Work{Id: Id(3 + (v * 26)), TaskId: Id(v + 1), Duration: 500, Resource: 8})
		worksBuffer = append(worksBuffer, Work{Id: Id(4 + (v * 26)), TaskId: Id(v + 1), Duration: 40000, Resource: 5, PreviousIds: third})
		worksBuffer = append(worksBuffer, Work{Id: Id(5 + (v * 26)), TaskId: Id(v + 1), Duration: 45000, Resource: 7, PreviousIds: firstAndThird})
		worksBuffer = append(worksBuffer, Work{Id: Id(6 + (v * 26)), TaskId: Id(v + 1), Duration: 4000, Resource: 8, PreviousIds: first})
		worksBuffer = append(worksBuffer, Work{Id: Id(7 + (v * 26)), TaskId: Id(v + 1), Duration: 3000, Resource: 2, PreviousIds: secondThirdAndFourth})
		worksBuffer = append(worksBuffer, Work{Id: Id(8 + (v * 26)), TaskId: Id(v + 1), Duration: 1000, Resource: 9})

		newChan := make(chan bool)
		goDoneChans = append(goDoneChans, newChan)

	}

	allDoneChannel := make(chan bool)

	doneSlice := make([]bool, 0)
	for h := 0; h < amount; h++ {
		doneSlice = append(doneSlice, false)
	}

	for g := range goDoneChans {
		doneSliceElement := &doneSlice[g]
		go func() {
			durationBuffer, doneSliceBuffer := calculate(worksBuffer[(8*g):8*(g+1)], doneSliceElement)
			result = append(result, resultStruct{Duration: durationBuffer, WorksPath: doneSliceBuffer})
			fmt.Println("another goroutine worked out...")
		}()
	}

	go func() {
		for {
			channelsDoneBoolSlice := make([]bool, amount)
			isDone := false
			for i := range channelsDoneBoolSlice {
				channelsDoneBoolSlice[i] = doneSlice[i]
				//fmt.Println(channelsDoneBoolSlice)
				if channelsDoneBoolSlice[i] {
					if i == len(channelsDoneBoolSlice)-1 {
						isDone = true
						fmt.Println(isDone, "<- done")
					}
				} else {
					break
				}
			}

			if isDone {
				allDoneChannel <- true
				break
			}
		}

	}()

	<-allDoneChannel
	fmt.Println("все сработало")

	allTasks = make([]Task, 0)
	worksBuffer = make([]Work, 0)
	goDoneChans = make([]chan bool, 0)

	resDuration := make([]int, 0)
	for x := range result {
		resDuration = append(resDuration, result[x].Duration)
	}

	res, ind := findMinValueAndIndex(resDuration)

	for a := range result[ind].WorksPath {
		result[ind].WorksPath[a].Id %= 26
		for b := range result[ind].WorksPath[a].PreviousIds {
			result[ind].WorksPath[a].PreviousIds[b] %= 26
		}
	}

	path := fmt.Sprintf("%v", result[ind].WorksPath)

	rdb := RedisConnect()
	RedisSet(rdb, "duration", strconv.Itoa(res))
	RedisSet(rdb, "path", path)
	fmt.Println(RedisGet(rdb, "path"), " <- это из редиса")

	return res, path, nil
}

func calculate(allWorksOfTaskSlice []Work, done *bool) (int, []Work) {
	defer func() { *done = true }()
	rand.Seed(time.Now().Unix())

	doneWorksSlice := make([]Work, 0)

	undoneWorksSlice := make([]Work, 0)

	availableWorksSlice := make([]Work, 0)

	inProgressWorksSlice := make([]Work, 0)

	resources := 10

	toDropAvailableIds := make([]Id, 0)

	toDropInProgressIds := make([]Id, 0)

	choosen := make([]Work, 0)

	var minDuration int

	var result int

	for _, workOfTask := range allWorksOfTaskSlice {
		undoneWorksSlice = append(undoneWorksSlice, workOfTask)
	}
	var u int
	for {
		u += 1
		fmt.Println(u)
		for _, workOfTask := range allWorksOfTaskSlice {
			//если ворка нет среди выполненных, доступных для выполнения или находящихся в обработке - проверяем дальше
			if (findWorkById(doneWorksSlice, workOfTask.Id) == -1) && (findWorkById(availableWorksSlice, workOfTask.Id) == -1) && (findWorkById(inProgressWorksSlice, workOfTask.Id) == -1) {
				var doneForTheWork int
				for _, prevId := range workOfTask.PreviousIds {
					//если работа из списка предшествующих есть среди выполненных, увеличиваем счетчик
					if findWorkById(doneWorksSlice, prevId) != -1 {
						doneForTheWork += 1
					}
				}
				if len(workOfTask.PreviousIds) == doneForTheWork { //если значение счетчика равно количеству предшествующих работ, значит работу можно выполнять - добавляем ее в available
					availableWorksSlice = append(availableWorksSlice, workOfTask)
				}
			}
		}

		//fmt.Println(availableWorksSlice, "available/")

		//сортируем слайс available по возрастанию
		sort.Slice(availableWorksSlice, func(i, j int) bool {
			return availableWorksSlice[i].Duration < availableWorksSlice[j].Duration
		})

		//fmt.Println(availableWorksSlice, "/available")

		var ran int
		for _ = range availableWorksSlice {
			randomAvailable := availableWorksSlice[rand.Intn(len(availableWorksSlice))]
			for findWorkById(choosen, randomAvailable.Id) != -1 {
				ran = rand.Intn(len(availableWorksSlice))
				randomAvailable = availableWorksSlice[ran]
				//fmt.Println(ran, "randomness")
			}

			//если позволяют ресурсы - добавляем работу в слайс "в обработке", уменьшаем доступные ресурсы на ресурсы работы и добавляем id в слайс для последующего удаления
			if resources >= randomAvailable.Resource {
				resources -= randomAvailable.Resource
				toDropAvailableIds = append(toDropAvailableIds, randomAvailable.Id)
				inProgressWorksSlice = append(inProgressWorksSlice, randomAvailable)
				choosen = append(choosen, randomAvailable)
				//fmt.Println(ran, "random")
			}
		}
		//удаляем элементы, добавленные в обработку из доступных
		for _, element := range toDropAvailableIds {
			availableWorksSlice = dropElementFromWorkSliceById(availableWorksSlice, element)
		}
		//очищаем слайс id
		toDropAvailableIds = make([]Id, 0)

		//находим минимальную длительность среди выполняемых работ, а также индекс соответствующего элемента
		_, minDuration = findMinDurationAndIndex(inProgressWorksSlice)
		//fmt.Println(inProgressWorksSlice, "pam")

		//добавляем найденную продолжительность в результат
		result += minDuration

		//fmt.Println(minDuration, "duration")

		//fmt.Println(inProgressWorksSlice, "in progress")

		//убавляем длительность всех работ на найденную минимальную
		for i := range inProgressWorksSlice {
			inProgressWorksSlice[i].Duration -= minDuration
			//если продолжительность == 0 - работа выполнена и ее надо удалить
			if inProgressWorksSlice[i].Duration == 0 {
				toDropInProgressIds = append(toDropInProgressIds, inProgressWorksSlice[i].Id)
				ind := findWorkById(allWorksOfTaskSlice, inProgressWorksSlice[i].Id)
				doneWork := allWorksOfTaskSlice[ind]
				doneWorksSlice = append(doneWorksSlice, doneWork)
				resources += inProgressWorksSlice[i].Resource
			}
		}

		//fmt.Println(toDropInProgressIds, "->")

		//удаляем выполненные работы из слайсов "в обработке" и undone
		for _, currentId := range toDropInProgressIds {
			inProgressWorksSlice = dropElementFromWorkSliceById(inProgressWorksSlice, currentId)

			undoneWorksSlice = dropElementFromWorkSliceById(undoneWorksSlice, currentId)
		}

		//fmt.Println(undoneWorksSlice, "<-")

		//очищаем слайс удаляемых из обработки работ
		toDropInProgressIds = make([]Id, 0)

		//когда нету невыполненных работ - возвращаем вычисленное значение
		if len(undoneWorksSlice) == 0 {
			fmt.Println(result, "итог")

			return result, doneWorksSlice
		}
	}
}

func findMinDurationAndIndex(slice []Work) (int, int) {
	min := 0
	index := -1

	fmt.Println(slice, "slice")

	for i, value := range slice {
		if value.Duration < min || i == 0 {
			min = value.Duration
			index = i
		}
	}

	fmt.Println(min, "min")
	return index, min
}

func findWorkById(slice []Work, id Id) int {
	for i, element := range slice {
		if element.Id == id {
			return i
		}
	}
	return -1
}

func dropElementFromWorkSliceById(slice []Work, id Id) []Work {
	index := findWorkById(slice, id)
	if index != -1 {
		if index != len(slice)-1 {
			return append(slice[:index], slice[index+1:]...)
		}
		return slice[:index]
	}
	return make([]Work, 0)
}
