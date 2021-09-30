package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/borud/broker"
	"github.com/ebobo/cpp_microservice_go/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var progress int = 0

func startProgess() {
	tk := time.NewTicker(1 * time.Second)
	for range tk.C {
		progress += 10
		// log.Printf("increase progress %d", progress)
		if progress >= 100 {
			tk.Stop()
			break
		}
	}
}

func reportProgess(conn *websocket.Conn) {
	tk := time.NewTicker(2 * time.Second)
	for range tk.C {
		// log.Printf("report message %d", progress)
		conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(progress)))
		if progress >= 100 {
			progress = 0
			conn.Close()
			break
		}
	}
}

func reportAnswer(conn *websocket.Conn, sub *broker.Subscriber) {
	defer func() {
		sub.Cancel()
		conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(60*time.Second))
		conn.Close()
	}()
	// log.Printf("report message %d", progress)
	msg := <-sub.Messages()
	json, err := json.Marshal(msg.Payload)

	if err != nil {
		log.Printf("error marshalling message to JSON: %v", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte(json))

}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello, This is calculator server`))
}

func (s *Server) setParameters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "get called"}`))
	case "POST":
		w.WriteHeader(http.StatusAccepted)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var para model.Parameter
		json.Unmarshal(reqBody, &para)
		json.NewEncoder(w).Encode(para)
		// log.Println(para)
		s.service.CalcService.Notify("1", &para)

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

func (s *Server) wsServices(w http.ResponseWriter, r *http.Request) {
	log.Println("ws services")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}
	log.Println("Web Socket Client successfully connected")

	sub, err := s.service.Broker.Subscribe("question_answered")
	if err != nil {
		log.Printf("subscription to '%s' failed: %v", "question_answered", err)
		http.Error(w, "subscription failed", http.StatusInternalServerError)
		return
	}

	go reportAnswer(ws, sub)

}

func (s *Server) startREST() {
	m := mux.NewRouter()
	m.HandleFunc("/api/parameters", s.setParameters).Methods("POST")
	m.HandleFunc("/api/ws", s.wsServices).Methods("GET")
	m.HandleFunc("/", home).Methods("GET")

	// Add CORS
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		MaxAge:           31,
		Debug:            false,
	})

	httpServer := &http.Server{
		Addr:              s.restAddr,
		Handler:           handlers.ProxyHeaders(cors.Handler(m)),
		ReadTimeout:       (10 * time.Second),
		ReadHeaderTimeout: (8 * time.Second),
		WriteTimeout:      (45 * time.Second),
	}

	// Shut down webserver when done
	go func() {
		<-s.ctx.Done()
		log.Printf("Shutting down REST interface")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Printf("Error shutting down REST interface: %v", err)
		}
		s.restStopped.Done()
	}()

	// Start webserver
	go func() {
		log.Printf("Starting REST at '%s'", s.restAddr)
		s.restStarted.Done()
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("REST interface: %v", err)
		} else {
			log.Printf("REST interface shut down")
		}
	}()
}
