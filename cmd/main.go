package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/handler"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/repository"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/service"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/server"

	"github.com/julienschmidt/httprouter"
)

func main() {

	db := repository.NewDb()

	wr := repository.NewWork(db)
	ws := service.NewWork(wr)
	wh := handler.NewWork(ws)

	r := httprouter.New()

	r.POST("/work", wrapHandler(wh.Create))

	theServer := server.New("8000", r)

	var err error

	go func() {
		err = theServer.Run()
	}()

	_ = repository.NewDb()

	fmt.Println("Server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err != nil {
		log.Fatalf("Error occured while running server - %s", err.Error())
	}
}

func wrapHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
