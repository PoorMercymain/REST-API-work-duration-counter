package main

import (
	api "github.com/PoorMercymain/REST-API-work-duration-counter/pkg/api/api/proto"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/counter"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &counter.GRPCServer{}
	api.RegisterCounterServer(s, srv)

	l, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
