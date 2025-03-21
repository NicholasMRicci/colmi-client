package lib

import (
	_ "embed"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/NicholasMRicci/colmi-client/lib/message"
	"github.com/gorilla/websocket"
	"tinygo.org/x/bluetooth"
)

type Server struct {
	mu       sync.Mutex
	bpm      map[int64]chan uint8
	stopRead chan struct{}
	mux      *http.ServeMux
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer() *Server {
	mux := http.NewServeMux()
	server := Server{mux: mux}

	mux.HandleFunc("/", ServePage)
	mux.HandleFunc("/bpm", server.ResgisterSocket())
	return &server
}

func (s *Server) Start() {
	go http.ListenAndServe("0.0.0.0:8080", s.mux)
}

func (s *Server) Stop() {
	if s.stopRead != nil {
		go close(s.stopRead)
	}
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	// Read index.html from ./index.html using file handling
	file, err := os.ReadFile("./lib/index.html")
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(file)
}

func (s *Server) ResgisterSocket() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		err = conn.WriteJSON("Connected")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("About to enter block 1")
		nonce := rand.Int63()
		bpmChan := make(chan byte, 1)
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			log.Println("Entering block 1")
			if len(s.bpm) == 0 {
				s.bpm = make(map[int64]chan uint8)
				go s.goGoRing()
			}

			s.bpm[nonce] = bpmChan
		}()

		defer func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.bpm, nonce)
			if len(s.bpm) == 0 {
				s.stopRead <- struct{}{}
			}
		}()
		// go func() {
		// 	for {
		// 		if _, _, err := conn.NextReader(); err != nil {
		// 			log.Println("The thing happened")
		// 			close(s.stopRead)
		// 			conn.Close()
		// 			break
		// 		}
		// 	}
		// }()
		for {
			log.Println("Waiting for bpm")
			bpm, ok := <-bpmChan
			if !ok {
				break
			}
			var err error
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if bpm == 0 {
				err = conn.WriteJSON("Measuring...")
			} else {
				err = conn.WriteJSON(bpm)
			}
			if err != nil {
				break
			}
		}
	}
}

func (s *Server) goGoRing() {
	// Start scanning.
	ring, err := AquireRing(bluetooth.DefaultAdapter, "")
	Must(err)
	Must(ring.Send(message.BlinkTwice()))

	s.stopRead = make(chan struct{}, 0)
	messages := make(chan message.Message)
	ring.BeginReads(messages)
	Must(ring.Send(message.StartWorkout()))
	defer func() {
		ring.StopReads()
		Must(ring.Send(message.PauseWorkout()))
		Must(ring.Send(message.EndWorkout()))
		Must(ring.Disconnect())
	}()
	for {
		select {
		case msg := <-messages:
			bpm, ok := message.DecodeWorkout(msg)
			if !ok {
				log.Printf("Unexpected msg: %v", msg)
			}
			s.mu.Lock()
			for _, val := range s.bpm {
				val <- bpm
			}
			s.mu.Unlock()
		case <-s.stopRead:
			log.Println("Got the stop signal")
			s.mu.Lock()
			defer s.mu.Unlock()
			for _, val := range s.bpm {
				close(val)
			}
			s.bpm = make(map[int64]chan uint8)
			return
		}
	}
}
