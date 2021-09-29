package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebobo/cpp_microservice_go/server"
)

type serverCmd struct {
	//short argu can only be single letter
	GRPCAddress string `short:"g" long:"grpc-addr" default:":5005" description:"gRPC listen address"`
	RESTAddress string `short:"r" long:"rest-addr" default:":5006" description:"REST listen address"`
}

// Execute export
func (sc *serverCmd) Execute(_ []string) error {
	fmt.Println("GRPC server at", sc.GRPCAddress)
	fmt.Println("REST server at", sc.RESTAddress)

	server := server.New(server.Config{
		GRPCAddr: sc.RESTAddress,
		RESTAddr: sc.GRPCAddress,
	})

	// Fire up the server
	err := server.Start()
	if err != nil {
		log.Fatalf("Unable to start srv: %v", err)
	}

	// Capture Ctrl-C from user
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	server.Shutdown()
	return nil
}
