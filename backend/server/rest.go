package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Parameter struct {
	A int32 `json:"A"`
	B int32 `json:"B"`
}

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

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello, This is calculator server`))
}

func setParameters(w http.ResponseWriter, r *http.Request) {
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
		go startProgess()

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

func wsServices(w http.ResponseWriter, r *http.Request) {
	log.Println("ws services")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}
	log.Println("Web Socket Client successfully connected")
	go reportProgess(ws)
}

func (s *Server) startREST() {
	m := mux.NewRouter()
	m.HandleFunc("/api/parameters", setParameters).Methods("POST")
	m.HandleFunc("/api/ws", wsServices).Methods("GET")
	m.HandleFunc("/", home).Methods("GET")

	// Add CORS
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST, GET, OPTIONS, PUT, DELETE"},
		MaxAge:           31,
		Debug:            true,
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
