package server

import (
	"context"
	"log"
	"net"

	"github.com/ebobo/cpp_microservice_go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClaculatorServer struct {
}

func NewClaculatorServer() *ClaculatorServer {
	return &ClaculatorServer{}
}

func (cs *ClaculatorServer) Questions(_ *emptypb.Empty, stream protos.Claculator_QuestionsServer) error {
	return nil
}

func (cs *ClaculatorServer) QuestionAnswered(ctx context.Context, in *protos.Answer) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

// startGRPC fires up the gRPC interface
func (s *Server) startGRPC() {
	clacServer := NewClaculatorServer()

	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		log.Fatalf("Error creating listen socket: %v", err)
	}

	gs := grpc.NewServer()

	// grpc display help
	reflection.Register(gs)

	protos.RegisterClaculatorServer(gs, clacServer)

	// Set up graceful shutdown
	go func() {
		<-s.ctx.Done()
		log.Printf("Shutting down gRPC interface")
		gs.GracefulStop()
	}()

	// Start gRPC server
	go func() {
		log.Printf("Starting gRPC at '%s'", s.grpcAddr)
		s.grpcStarted.Done()
		err = gs.Serve(listener)
		if err != nil {
			log.Printf("gRPC interface returned error: %v", err)
		}
		log.Printf("gRPC interface: shut down")
		s.grpcStopped.Done()
	}()

}
