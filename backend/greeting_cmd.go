package main

import "fmt"

type greetingCmd struct {
	Name string `short:"n" long:"name" description:"Your name, for a greeting" default:"Unknown"`
}

// Execute export
func (r *greetingCmd) Execute(args []string) error {
	fmt.Println("Hello", r.Name)
	return nil
}
