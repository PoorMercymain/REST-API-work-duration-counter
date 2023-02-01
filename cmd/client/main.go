package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/handler"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/repository"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/service"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/router"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/server"

	"github.com/julienschmidt/httprouter"
)

func main() {

	fmt.Println(build.Default.GOPATH, "GOPATH")
	rdb := repository.RedisConnect()
	repository.RedisSet(rdb, "it", "works!")
	fmt.Println(repository.RedisGet(rdb, "it"), "HERE!!!")

	db := repository.NewDb()

	wr := repository.NewWork(db)
	ws := service.NewWork(wr)
	wh := handler.NewWork(ws)

	tr := repository.NewTask(db)
	ts := service.NewTask(tr)
	th := handler.NewTask(ts)

	r := httprouter.New()

	r.POST("/work", router.WrapHandler(wh.Create))
	r.DELETE("/work/:id", router.WrapHandler(wh.Delete))
	r.GET("/works/:id", router.WrapHandler(wh.List))
	r.GET("/createTestWorks/", router.WrapHandler(wh.CreateTestWorks))

	r.POST("/task", router.WrapHandler(th.Create))
	r.DELETE("/task/:id", router.WrapHandler(th.Delete))
	r.PUT("/task", router.WrapHandler(th.Update))
	r.GET("/taskWorks/:id", router.WrapHandler(th.ListWorksOfTask))
	r.GET("/taskCount/:id", router.WrapHandler(th.CountDuration))
	r.GET("/createTestTasks/", router.WrapHandler(th.CreateTestTasks))
	r.GET("/countAll/", router.WrapHandler(th.CountAll))

	theServer := server.New("8000", r)

	var err error

	go func() {
		err = theServer.Run()
	}()

	fmt.Println("Server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err != nil {
		log.Fatalf("Error occured while running server - %s", err.Error())
	}
}
