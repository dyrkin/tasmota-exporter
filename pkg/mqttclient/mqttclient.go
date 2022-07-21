package mqttclient

import (
	"fmt"
	"log"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	c mqtt.Client
}

func NewMqttClient(host string, port int, user, password string) *MqttClient {
	mqttClient := &MqttClient{}
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	options.SetClientID("prometheus_tasmota_exporter")
	options.SetUsername(user)
	options.SetPassword(password)
	options.SetCleanSession(false)
	options.SetAutoReconnect(true)
	options.OnConnect = mqttClient.connectionHandler
	options.OnConnectionLost = mqttClient.connectionLostHandler

	mqttClient.c = mqtt.NewClient(options)

	return mqttClient
}

func (mc *MqttClient) Connect() error {
	token := mc.c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mc *MqttClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	token := mc.c.Subscribe(topic, 1, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mc *MqttClient) SendCommand(topic, payload string) error {
	token := mc.c.Publish(topic, 1, false, []byte(payload))
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mc *MqttClient) connectionHandler(_ mqtt.Client) {
	log.Printf("mqtt connected")
}

func (mc *MqttClient) connectionLostHandler(_ mqtt.Client, err error) {
	log.Printf("mqtt disconnected. reason: %s", err)
	log.Println("exiting...")
	syscall.Exit(1)
}
