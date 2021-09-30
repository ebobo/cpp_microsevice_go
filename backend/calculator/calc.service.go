package main

import (
	"fmt"
	"io"
	"log"
	"time"
)

type CalcMicroService struct {
	grpc_client *GrpcClient
}

func NewCalcMicroService() *CalcMicroService {
	return &CalcMicroService{}
}

// Run runs whole algorithm to process maps
func (cms *CalcMicroService) Run() {
	log.Println("Running calculation micro service")
	cms.grpc_client = NewGrpcClient("localhost:5006")

	err := cms.grpc_client.Init()
	if err != nil {
		log.Println(err)
		return
	}

	done := make(chan bool)

	go func() {
		for {
			msg, err := cms.grpc_client.GetQuestions()
			if err == io.EOF {
				fmt.Println("EOF")
				// FIXME: reconnect
				done <- true //means stream is finished
				return
			}
			fmt.Println("question comming in", msg)
			id, result, err := cms.doCalculation(msg.Id, msg.A, msg.B, msg.Type)
			fmt.Printf("Result=%v \n", result)
			if err != nil {
				fmt.Println(err)
			} else {
				cms.grpc_client.SendAnswer(id, result)
			}
		}
	}()
	<-done //we will wait until all response is received
}

func (cms *CalcMicroService) doCalculation(id string, pa int32, pb int32, typ string) (string, int32, error) {
	time.Sleep(2 * time.Second)
	var res int32 = 0
	switch typ {
	case "plus":
		res = pa + pb
	case "minus":
		res = pa - pb
	default:
		res = 0
	}

	return id, res, nil
}
