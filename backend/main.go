package main

import (
	"log"
	"net/http"

	"github.com/ebobo/cpp_microservice_go/cors"
)

type Parameter struct {
    A int32 `json:"A"`
	B int32 `json:"B"`
	Type string `json:"type"`
}

func setParameters (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "get called"}`))
    case "POST":
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "post called"}`))
    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "put called"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "delete called"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}

func main() {
	log.Println("server is running on port 5006")
	
	parameterHandler := http.HandlerFunc(setParameters)
    http.Handle("/parameters", cors.Middleware(parameterHandler))

	err := http.ListenAndServe(":5006", nil)
	if err != nil {
		log.Fatal(err)
	}
}