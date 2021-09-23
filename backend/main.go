package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ebobo/cpp_microservice_go/cors"
	"github.com/gorilla/websocket"
)

type Parameter struct {
    A int32 `json:"A"`
	B int32 `json:"B"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var progress = 0


func startProgess() {
  	tk := time.NewTicker(1 * time.Second)
	for range tk.C {
		progress += 10
		log.Printf("increase progress %d", progress)  
		if progress >= 100 {
			tk.Stop()
			break
		}
	}
}

func reportProgess(conn *websocket.Conn) {
	tk := time.NewTicker(2 * time.Second)
  	for range tk.C {
		log.Printf("report message %d", progress)  
		conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(progress)))

	  if progress >= 100 {
		  progress = 0;
		  conn.Close()
		  break
	  }
  }
}

func setParameters (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "get called"}`))
    case "POST":
		w.WriteHeader(http.StatusAccepted)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var para Parameter 
		json.Unmarshal(reqBody, &para)
		json.NewEncoder(w).Encode(para)
		
		go startProgess();
		
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


func wsServices (w http.ResponseWriter, r *http.Request) {
	log.Println("ws services")
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
	ws, err :=upgrader.Upgrade(w,r,nil)

	if err != nil {
		log.Println(err)
	}

    log.Println("Web Socket Client successfully connected")
  		
	go reportProgess(ws)		
}


func main() {
	log.Println("server is running on port 5006")

	paraHandler := http.HandlerFunc(setParameters)
	http.Handle("/api/parameters", cors.Middleware(paraHandler))

	wsHandler := http.HandlerFunc(wsServices)
	http.Handle("/api/ws", cors.Middleware(wsHandler))

	err := http.ListenAndServe(":5006", nil)
	if err != nil {
		log.Fatal(err)
	}
}