package service

import (
	"context"
	"sync"

	"github.com/ebobo/cpp_microservice_go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Service represents the service that serves both the REST and gRPC
// APIs.  The service is what plugs the API into the various aspects
// of the backend.  Mostly persistence for maintaining configuration.

// In Go, an identifier that starts with a capital letter is exported \
// from the package, and can be accessed by anyone outside the package
// that declares it.
type Service struct {
	CalcServer         *ClaculatorServer
	started            *sync.WaitGroup
	stopped            *sync.WaitGroup
	ctx                context.Context
	cancel             context.CancelFunc
	internalClientConn *grpc.ClientConn
}

func NewService() *Service {
	return &Service{}
}

// GRPCServer returns a grpc.Server instance with the services
// populated.
func (s *Service) GRPCServer() *grpc.Server {
	gs := grpc.NewServer()
	s.CalcServer = NewClaculatorServer()

	reflection.Register(gs)

	protos.RegisterClaculatorServer(gs, s.CalcServer)

	return gs
}
