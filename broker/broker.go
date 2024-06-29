package main

import (
	"context"
	"log"
	"net"

	pb "github.com/felipefferrada/lab5/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBrokerServer
}

func (s *server) RedirigirComando(ctx context.Context, req *pb.ComandoRequest) (*pb.ComandoResponse, error) {
	// Implementación de lógica para redirigir comando
	return &pb.ComandoResponse{DireccionFulcrum: "localhost:5001"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &server{})

	log.Printf("Broker server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
