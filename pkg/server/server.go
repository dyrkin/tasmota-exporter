package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dyrkin/tasmota-exporter/pkg/metrics"
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

func (s *Server) Start() error {
	slog.Info("Started listening", "port", s.port)
	return s.httpServer.ListenAndServe()
}
