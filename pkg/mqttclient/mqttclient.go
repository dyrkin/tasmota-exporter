package mqttclient

import (
	"fmt"
	"log/slog"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	c mqtt.Client
}

func NewMqttClient(host string, port int, user, password string, client_id string) *MqttClient {
	mqttClient := &MqttClient{}
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	options.SetClientID(client_id)
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

func (mc *MqttClient) connectionHandler(c mqtt.Client) {
	slog.Info("mqtt connected")
}

func (mc *MqttClient) connectionLostHandler(_ mqtt.Client, err error) {
	slog.Error("mqtt disconnected, exiting.", "reason", err)
	syscall.Exit(1)
}
