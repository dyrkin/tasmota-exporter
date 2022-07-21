package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dyrkin/tasmota-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	port       int
}

func NewServer(port int, metrics *metrics.Metrics) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	s := &Server{
		httpServer: httpServer,
		port:       port,
	}

	mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		metrics.Refresh()
		promhttp.Handler().ServeHTTP(writer, request)
	})

	return s
}

func (s *Server) Start() {
	log.Printf("Started listening on: %d", s.port)

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
