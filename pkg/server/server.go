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
}

func NewServer(port int, metrics *metrics.Metrics) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	s := &Server{
		httpServer: httpServer,
	}

	mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		metrics.Refresh()
		promhttp.Handler().ServeHTTP(writer, request)
	})

	return s
}

func (s *Server) Start() {
	log.Println("Starting server")

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
