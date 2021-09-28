package main

import (
	"context"
	"log"

	"github.com/ebobo/cpp_microservice_go/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GrpcClient and connect by gRPC
type GrpcClient struct {
	serverURL string
	client    protos.ClaculatorClient
}

func NewGrpcClient(url string) *GrpcClient {
	return &GrpcClient{serverURL: url}
}

// Init initializes grpc client, connects to the server
func (c *GrpcClient) Init() error {
	conn, err := grpc.Dial(c.serverURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
		return err
	}
	c.client = protos.NewClaculatorClient(conn)
	return nil
}

//  creates stream and wait until receive message
//  on the stream about question raised
func (c *GrpcClient) GetQuestions() (*protos.QuestionRaised, error) {
	stream, err := c.client.Questions(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("open stream error  %v", err)
		return nil, err
	}
	return stream.Recv()
}

// Send Answer to grpc server
func (c *GrpcClient) SendAnswer(id string, res int32) error {
	_, err := c.client.QuestionAnswered(context.Background(), &protos.Answer{Id: id, Result: res})
	return err
}
