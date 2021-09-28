package main

type serverCmd struct {
	//short argu can only be single letter
	GRPCAddress string `short:"g" long:"grpc-addr" default:":5005" description:"gRPC listen address"`
	RESTAddress string `short:"r" long:"rest-addr" default:":5006" description:"REST listen address"`
}

// Execute export
func (sc *serverCmd) Execute(args []string) error {
	// fmt.Println("GRPC", sc.GRPCAddress)
	// fmt.Println("REST", sc.RESTAddress)

	return nil
}
