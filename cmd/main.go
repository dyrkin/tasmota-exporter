package main

import (
	"log/slog"
	"os"

	"github.com/dyrkin/tasmota-exporter/pkg/engine"
	"github.com/dyrkin/tasmota-exporter/pkg/metrics"

	"github.com/dyrkin/tasmota-exporter/pkg/mqttclient"
	"github.com/dyrkin/tasmota-exporter/pkg/server"
)

func abort(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func main() {
	v, err := ReadEnv()
	if err != nil {
		abort("Can't read env variables, exiting.", "error", err)
	}
	pm := metrics.NewPlainMetrics()
	m := metrics.NewMetrics(pm)
	c := metrics.NewCleaner(m, v.removeWhenInactiveMinutes)
	c.Start()

	mqttClient := mqttclient.NewMqttClient(v.mqttHost, v.mqttPort, v.mqttUsername, v.mqttPassword, v.mqttClientId)
	if err := mqttClient.Connect(); err != nil {
		abort("Can't connect to mqtt broker", "error", err)
	}

	e := engine.NewEngine(mqttClient, pm, v.statusUpdateSeconds)
	if err := e.Subscribe(v.mqttTopics); err != nil {
		abort("Can't subscribe to topics", "error", err)
	}

	s := server.NewServer(v.serverPort, m)
	if err := s.Start(); err != nil {
		abort("Failed to start server", "error", err)
	}
}
