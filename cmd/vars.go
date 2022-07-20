package main

import (
	"log"
	"os"
	"strconv"
)

type vars struct {
	mqttHost, mqttUsername, mqttPassword string
	mqttPort, serverPort                 int
}

func ReadEnv() *vars {
	v := &vars{}
	v.mqttHost = orDefault(os.Getenv("MQTT_HOSTNAME"), "localhost")
	mqttPort, err := strconv.Atoi(orDefault(os.Getenv("MQTT_PORT"), "1883"))
	if err != nil {
		log.Fatalf("can't parse provided mqtt port: %s", err)
	}
	v.mqttPort = mqttPort
	v.mqttUsername = orDefault(os.Getenv("MQTT_USERNAME"), "")
	v.mqttPassword = orDefault(os.Getenv("MQTT_PASSWORD"), "")
	serverPort, err := strconv.Atoi(orDefault(os.Getenv("PROMETHEUS_EXPORTER_PORT"), "9092"))
	if err != nil {
		log.Fatalf("can't parse provided server port: %s", err)
	}
	v.serverPort = serverPort
	return v
}
