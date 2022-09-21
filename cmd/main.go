package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"server/pkg/server"
	"syscall"
)

func main() {
	theServer := server.New("8000")

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
