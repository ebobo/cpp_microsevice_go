package main

import (
	"log"
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
}
