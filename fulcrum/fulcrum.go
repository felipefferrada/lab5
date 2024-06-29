package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/felipefferrada/lab5/proto"
	"google.golang.org/grpc"
)

type FulcrumServer struct {
	pb.UnimplementedFulcrumServer
	mu       sync.Mutex
	sectores map[string]map[string]int32 // Estructura para almacenar datos
}

func (s *FulcrumServer) AgregarBase(ctx context.Context, req *pb.AgregarBaseRequest) (*pb.AgregarBaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sectores == nil {
		s.sectores = make(map[string]map[string]int32)
	}

	if _, ok := s.sectores[req.Sector]; !ok {
		s.sectores[req.Sector] = make(map[string]int32)
	}

	s.sectores[req.Sector][req.Base] = req.Enemigos

	return &pb.AgregarBaseResponse{Mensaje: "Base agregada exitosamente"}, nil
}

func (s *FulcrumServer) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.RenombrarBaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sectores == nil {
		return &pb.RenombrarBaseResponse{Mensaje: "No se encontraron sectores"}, nil
	}

	if bases, ok := s.sectores[req.Sector]; ok {
		if enemigos, ok := bases[req.BaseAntigua]; ok {
			delete(bases, req.BaseAntigua)
			bases[req.BaseNueva] = enemigos
			return &pb.RenombrarBaseResponse{Mensaje: "Base renombrada exitosamente"}, nil
		}
	}

	return &pb.RenombrarBaseResponse{Mensaje: "Base no encontrada"}, nil
}

func (s *FulcrumServer) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.ActualizarValorResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sectores == nil {
		return &pb.ActualizarValorResponse{Mensaje: "No se encontraron sectores"}, nil
	}

	if bases, ok := s.sectores[req.Sector]; ok {
		if _, ok := bases[req.Base]; ok {
			bases[req.Base] = req.Enemigos
			return &pb.ActualizarValorResponse{Mensaje: "Valor actualizado exitosamente"}, nil
		}
	}

	return &pb.ActualizarValorResponse{Mensaje: "Base no encontrada"}, nil
}

func (s *FulcrumServer) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.BorrarBaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sectores == nil {
		return &pb.BorrarBaseResponse{Mensaje: "No se encontraron sectores"}, nil
	}

	if bases, ok := s.sectores[req.Sector]; ok {
		if _, ok := bases[req.Base]; ok {
			delete(bases, req.Base)
			return &pb.BorrarBaseResponse{Mensaje: "Base borrada exitosamente"}, nil
		}
	}

	return &pb.BorrarBaseResponse{Mensaje: "Base no encontrada"}, nil
}

func (s *FulcrumServer) GetEnemigos(ctx context.Context, req *pb.GetEnemigosRequest) (*pb.GetEnemigosResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sectores == nil {
		return &pb.GetEnemigosResponse{Enemigos: 0}, nil
	}

	if bases, ok := s.sectores[req.Sector]; ok {
		if enemigos, ok := bases[req.Base]; ok {
			return &pb.GetEnemigosResponse{Enemigos: enemigos}, nil
		}
	}

	return &pb.GetEnemigosResponse{Enemigos: 0}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFulcrumServer(s, &FulcrumServer{})

	log.Printf("Fulcrum server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
