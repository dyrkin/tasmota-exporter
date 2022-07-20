package main

import (
	"github.com/dyrkin/tasmota-exporter/pkg/engine"
	"github.com/dyrkin/tasmota-exporter/pkg/metrics"
	"log"

	"github.com/dyrkin/tasmota-exporter/pkg/mqttclient"
	"github.com/dyrkin/tasmota-exporter/pkg/server"
)

func main() {
	v := ReadEnv()
	pm := metrics.NewPlainMetrics()
	m := metrics.NewMetrics(pm)
	mqttClient := mqttclient.NewMqttClient(v.mqttHost, v.mqttPort, v.mqttUsername, v.mqttPassword)
	if err := mqttClient.Connect(); err != nil {
		log.Fatalf("can't connect to mqtt broker: %s", err)
	}

	e := engine.NewEngine(mqttClient, pm)
	e.Subscribe([]string{"tele/+/+", "stat/+/+"})

	s := server.NewServer(v.serverPort, m)
	s.Start()
}

func orDefault(value string, def string) string {
	if value == "" {
		return def
	}
	return value
}
