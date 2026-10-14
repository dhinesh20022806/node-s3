package services

import (
	"net"

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

	// register the server


	return server.Serve(lis)
}