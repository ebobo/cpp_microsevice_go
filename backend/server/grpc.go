package server

import (
	"log"
	"net"
)

// startGRPC fires up the gRPC interface
func (s *Server) startGRPC() {
	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		log.Fatalf("Error creating listen socket: %v", err)
	}

	grpcServer := s.service.GRPCServer()

	// Set up graceful shutdown
	go func() {
		<-s.ctx.Done()
		log.Printf("Shutting down gRPC interface")
		grpcServer.GracefulStop()
	}()

	// Start gRPC server
	go func() {
		log.Printf("Starting gRPC at '%s'", s.grpcAddr)
		s.grpcStarted.Done()
		err = grpcServer.Serve(listener)
		log.Printf("Starting 2")

		if err != nil {
			log.Printf("gRPC interface returned error: %v", err)
		}
		log.Printf("gRPC interface: shut down")
		s.grpcStopped.Done()
	}()

}
