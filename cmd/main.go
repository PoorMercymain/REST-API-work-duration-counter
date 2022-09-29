package main

import (
	"fmt"
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

	db := repository.NewDb()

	wr := repository.NewWork(db)
	ws := service.NewWork(wr)
	wh := handler.NewWork(ws)

	tr := repository.NewTask(db)
	ts := service.NewTask(tr)
	th := handler.NewTask(ts)

	r := httprouter.New()

	//TODO: create task and work routes
	r.POST("/work", router.WrapHandler(wh.Create))
	r.POST("/task", router.WrapHandler(th.Create))

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
