package main

import (
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"google.golang.org/grpc"
	"log"
	"net"
	"shape/services"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	shape.RegisterShapeServer(server, services.Server{})

	log.Println("Started gRPC Server on port 3000")
	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
