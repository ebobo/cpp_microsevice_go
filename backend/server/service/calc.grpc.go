package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/borud/broker"
	"github.com/ebobo/cpp_microservice_go/model"
	"github.com/ebobo/cpp_microservice_go/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClaculatorService struct {
	questionStream protos.Claculator_QuestionsServer
	ctx            context.Context
	broker         *broker.Broker
}

func New(c context.Context, b *broker.Broker) *ClaculatorService {
	return &ClaculatorService{
		ctx:    c,
		broker: b,
	}
}

func (cs *ClaculatorService) Questions(_ *emptypb.Empty, stream protos.Claculator_QuestionsServer) error {
	cs.questionStream = stream
	<-cs.ctx.Done()
	return nil
}

func (cs *ClaculatorService) QuestionAnswered(ctx context.Context, in *protos.Answer) (*emptypb.Empty, error) {
	fmt.Println("QuestionAnswered ", in.Result)
	cs.broker.Publish("question_answered", in, time.Millisecond*5)
	return &emptypb.Empty{}, nil
}

// Notify the client for Question added
func (cs *ClaculatorService) Notify(id string, para *model.Parameter) {
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
