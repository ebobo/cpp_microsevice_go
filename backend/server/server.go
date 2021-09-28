package server

import (
	"context"
	"log"
	"sync"
)

// Server is the Calculation Backend Server
type Server struct {
	grpcAddr    string
	restAddr    string
	ctx         context.Context
	cancel      context.CancelFunc
	grpcStarted *sync.WaitGroup
	restStarted *sync.WaitGroup
	restStopped *sync.WaitGroup
	grpcStopped *sync.WaitGroup
}

// Config for Calculation Backend Server
type Config struct {
	GRPCAddr string
	RESTAddr string
}

// New returns a configured backend server
func New(config Config) *Server {
	return &Server{
		grpcAddr: config.GRPCAddr,
		restAddr: config.RESTAddr,
	}
}

// Start starts the server.
func (s *Server) Start() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.restStarted = &sync.WaitGroup{}
	s.restStopped = &sync.WaitGroup{}
	s.grpcStopped = &sync.WaitGroup{}
	s.grpcStarted = &sync.WaitGroup{}

	// Start gRPC interface
	// s.grpcStarted.Add(1)
	// s.grpcStopped.Add(1)
	// s.startGRPC()
	// s.grpcStarted.Wait()

	// REST interface
	if s.restAddr != "" {
		s.restStarted.Add(1)
		s.restStopped.Add(1)
		s.startREST()
		s.restStarted.Wait()
	} else {
		log.Printf("Skipping REST interface (no listen address given)")
	}

	return nil
}

// Shutdown takes down the network interfaces and stops the servers.
func (s *Server) Shutdown() {

	// Wait until interfaces shut down
	// s.grpcStopped.Wait()
	s.restStopped.Wait()
}
