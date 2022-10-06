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
	c := metrics.NewCleaner(m, v.removeWhenInactiveMinutes)
	c.Start()

	mqttClient := mqttclient.NewMqttClient(v.mqttHost, v.mqttPort, v.mqttUsername, v.mqttPassword)
	if err := mqttClient.Connect(); err != nil {
		log.Fatalf("can't connect to mqtt broker: %s", err)
	}

	e := engine.NewEngine(mqttClient, pm)
	e.Subscribe([]string{"tele/+/+", "stat/+/+"})

	s := server.NewServer(v.serverPort, m)
	s.Start()
}
