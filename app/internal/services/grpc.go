package services

import (
	"fmt"
	"net"

	"github.com/dhinesh20022806/node-s3/app/internal/services/object/handlers"
	"github.com/dhinesh20022806/node-s3/app/internal/services/object/service"
	"google.golang.org/grpc"
)



type grpcServer struct {
	address string
}

func NewGrpcServer(address string) *grpcServer {
	return &grpcServer{
		address: address,
	}
}

func (g *grpcServer) Run() error {
	lis, err := net.Listen("tcp", g.address)

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	objectServices := service.NewObject()

	handlers.NewObjectHandler(server, objectServices)

	// register the server
	fmt.Printf("grpc service were running: %s ", g.address)


	return server.Serve(lis)
}