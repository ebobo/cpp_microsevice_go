package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ebobo/cpp_microservice_go/cors"
)

type Parameter struct {
    A int32 `json:"A"`
	B int32 `json:"B"`
}


func setParameters (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "get called"}`))
    case "POST":
		reqBody, _ := ioutil.ReadAll(r.Body)
		var para Parameter 
		json.Unmarshal(reqBody, &para)
		log.Println(para)
        
		json.NewEncoder(w).Encode(para)
		
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

	paraHandler := http.HandlerFunc(setParameters)
	http.Handle("/api/parameters", cors.Middleware(paraHandler))

	err := http.ListenAndServe(":5006", nil)
	if err != nil {
		log.Fatal(err)
	}
}