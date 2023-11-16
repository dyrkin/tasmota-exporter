package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type vars struct {
	mqttHost, mqttUsername, mqttPassword, mqttClientId                   string
	mqttPort, serverPort, removeWhenInactiveMinutes, statusUpdateSeconds int
	mqttTopics                                                           []string
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
	v.mqttClientId = orDefault(os.Getenv("MQTT_CLIENT_ID"), "prometheus_tasmota_exporter")
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
	v.mqttTopics = orDefaultList(os.Getenv("MQTT_TOPICS"), "tele/+/+, stat/+/+")
	statusUpdateSeconds, err := strconv.Atoi(orDefault(os.Getenv("STATUS_UPDATE_SECONDS"), "5"))
	if err != nil {
		log.Fatalf("can't parse provided status update interval: %s", err)
	}
	v.statusUpdateSeconds = statusUpdateSeconds
	return v
}

func orDefault(value string, def string) string {
	if value == "" {
		return def
	}
	return value
}

// read env variable, split per comma, trim the values, and return them as a list of strings
func orDefaultList(value string, def string) []string {
	if value == "" {
		return splitAndTrim(def)
	}
	return splitAndTrim(value)
}

func splitAndTrim(value string) []string {
	s := strings.Split(value, ",")
	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}
	return s
}
