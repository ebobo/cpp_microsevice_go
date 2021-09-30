package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/borud/broker"
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
	CalcService *ClaculatorService
	Broker      *broker.Broker
	stopped     *sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewService(c context.Context, ca context.CancelFunc) *Service {
	var config = broker.Config{DownStreamChanLen: 0, PublishChanLen: 0, SubscribeChanLen: 0, UnsubscribeChanLen: 0, DeliveryTimeout: 0}

	return &Service{
		Broker: broker.New(config),
		ctx:    c,
		cancel: ca,
	}
}

// GRPCServer returns a grpc.Server instance with the services
// populated.
func (s *Service) GRPCServer() *grpc.Server {
	gs := grpc.NewServer()
	s.CalcService = New(s.ctx, s.Broker)

	reflection.Register(gs)

	protos.RegisterClaculatorServer(gs, s.CalcService)

	return gs
}

// Shutdown the service
func (s *Service) Shutdown() {
	fmt.Println("shut down service")
	if s.cancel != nil {
		s.cancel()
	}

}
