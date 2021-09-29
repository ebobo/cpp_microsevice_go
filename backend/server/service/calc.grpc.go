package service

import (
	"context"
	"log"

	"github.com/ebobo/cpp_microservice_go/model"
	"github.com/ebobo/cpp_microservice_go/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClaculatorServer struct {
	questionStream protos.Claculator_QuestionsServer
}

func NewClaculatorServer() *ClaculatorServer {
	return &ClaculatorServer{}
}

func (cs *ClaculatorServer) Questions(_ *emptypb.Empty, stream protos.Claculator_QuestionsServer) error {
	error := make(chan error)
	cs.questionStream = stream
	return <-error
}

func (cs *ClaculatorServer) QuestionAnswered(ctx context.Context, in *protos.Answer) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Notify the client for Question added
func (cs *ClaculatorServer) Notify(id string, para *model.Parameter) {
	m := protos.QuestionRaised{}
	m.Id = id
	m.A = para.A
	m.B = para.B
	m.Type = "plus"

	err := cs.questionStream.Send(&m)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

}
