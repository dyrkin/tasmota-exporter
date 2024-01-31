package engine

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/dyrkin/tasmota-exporter/pkg/metrics"
	"github.com/dyrkin/tasmota-exporter/pkg/mqttclient"
)

type Engine struct {
	scheduled           map[string]any
	lock                *sync.Mutex
	mqttClient          *mqttclient.MqttClient
	plainMetrics        *metrics.PlainMetrics
	statusUpdateSeconds int
}

func NewEngine(mqttClient *mqttclient.MqttClient, pm *metrics.PlainMetrics, statusUpdateSeconds int) *Engine {
	return &Engine{map[string]any{}, &sync.Mutex{}, mqttClient, pm, statusUpdateSeconds}
}

func (e *Engine) Subscribe(mqttListenTopics []string) error {
	for _, topic := range mqttListenTopics {
		if err := e.mqttClient.Subscribe(topic, e.messageProcessor); err != nil {
			return fmt.Errorf("can't subscribe to %q: %w", topic, err)
		}
	}
	return nil
}

func (e *Engine) messageProcessor(_ mqtt.Client, m mqtt.Message) {
	e.scheduleStatusCommand(m.Topic())
	rawMetrics := metrics.Extract(m.Payload())
	e.plainMetrics.Update(m.Topic(), rawMetrics)
}

func (e *Engine) scheduleStatusCommand(topic string) {
	if strings.HasPrefix(topic, "tele/") && strings.HasSuffix(topic, "/STATE") {
		segments := strings.Split(topic, "/")
		source := strings.Join(segments[1:len(segments)-1], "/")
		segments[0], segments[len(segments)-1] = "cmnd", "Status0"
		target := strings.Join(segments, "/")
		e.lock.Lock()
		defer e.lock.Unlock()
		if _, ok := e.scheduled[source]; !ok {
			slog.Debug("scheduling status updates", "interval_sec", e.statusUpdateSeconds, "source", source)
			e.scheduled[source] = true
			ticker := time.NewTicker(time.Duration(e.statusUpdateSeconds) * time.Second)
			go func() {
				for {
					select {
					case <-ticker.C:
						slog.Debug("sending status update request", "command", target)
						if err := e.mqttClient.SendCommand(target, ""); err != nil {
							slog.Error("can't send status command", "command", target)
						}
					}
				}
			}()
		}
	}
}
