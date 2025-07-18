package mysse

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	streams map[string]*stream
	lock    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		streams: make(map[string]*stream),
	}
}

func (s *Server) CreateStream(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.streams[name]; ok {
		return
	}
	st := &stream{
		name:       name,
		clients:    make(map[string]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan *Event, 16),
	}
	s.streams[name] = st
	go st.run()
}

func (s *Server) Publish(stream string, event *Event) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if st, ok := s.streams[stream]; ok {
		st.broadcast <- event
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		http.Error(w, "missing stream", http.StatusBadRequest)
		return
	}

	s.lock.RLock()
	st, ok := s.streams[streamName]
	s.lock.RUnlock()
	if !ok {
		http.Error(w, "stream not found", http.StatusNotFound)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	cl := &client{
		id:      fmt.Sprintf("%d", time.Now().UnixNano()),
		writer:  w,
		flusher: flusher,
		closeCh: make(chan struct{}),
		stream:  streamName,
	}
	st.register <- cl

	notify := r.Context().Done()
	<-notify
	st.unregister <- cl
}

// Internal types (unexported)
type client struct {
	id      string
	writer  http.ResponseWriter
	flusher http.Flusher
	closeCh chan struct{}
	stream  string
}

type stream struct {
	name       string
	clients    map[string]*client
	register   chan *client
	unregister chan *client
	broadcast  chan *Event
	lock       sync.RWMutex
}

func (s *stream) run() {
	for {
		select {
		case cl := <-s.register:
			s.lock.Lock()
			s.clients[cl.id] = cl
			s.lock.Unlock()

		case cl := <-s.unregister:
			s.lock.Lock()
			if _, ok := s.clients[cl.id]; ok {
				delete(s.clients, cl.id)
				close(cl.closeCh)
			}
			s.lock.Unlock()

		case event := <-s.broadcast:
			s.lock.RLock()
			for _, cl := range s.clients {
				_, _ = cl.writer.Write(event.Encode())
				cl.flusher.Flush()
			}
			s.lock.RUnlock()
		}
	}
}
