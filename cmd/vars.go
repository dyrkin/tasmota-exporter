package main

import (
	"log"
	"os"
	"strconv"
)

type vars struct {
	mqttHost, mqttUsername, mqttPassword            string
	mqttPort, serverPort, removeWhenInactiveMinutes int
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
	removeWhenInactiveMinutes, err := strconv.Atoi(orDefault(os.Getenv("REMOVE_WHEN_INACTIVE_MINUTES"), "1"))
	if err != nil {
		log.Fatalf("can't parse provided timeout: %s", err)
	}
	v.removeWhenInactiveMinutes = removeWhenInactiveMinutes
	return v
}

func orDefault(value string, def string) string {
	if value == "" {
		return def
	}
	return value
}
