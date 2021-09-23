package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Options contains the command line options
type Options struct {
	Message  string `short:"m" long:"message" env:"MESSAGE" description:"The input message" required:"yes"`
}

func main() {
	var opt Options
	parser := flags.NewParser(&opt, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	fmt.Printf("message %s\n", opt.Message)
}