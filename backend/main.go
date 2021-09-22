package main

import (
	"log"
	"net/http"
)

const basePath = "/api"

func main() {
	log.Println("server is running on port 5006")
	err := http.ListenAndServe(":5005", nil)
	if err != nil {
		log.Fatal(err)
	}
}