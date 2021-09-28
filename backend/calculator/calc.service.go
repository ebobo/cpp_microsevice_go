package main

import (
	"log"
)

type CalcMicroService struct {
	grpc_client *GrpcClient
}

// Run runs whole algorithm to process maps
func (cms *CalcMicroService) Run() {
	log.Println("Running calculation micro service")
	cms.grpc_client = NewGrpcClient("localhost:7601")

	err := cms.grpc_client.Init()
	if err != nil {
		log.Println(err)
		return
	}
}
