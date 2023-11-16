package engine

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dyrkin/tasmota-exporter/pkg/metrics"
	"github.com/dyrkin/tasmota-exporter/pkg/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Engine struct {
	scheduled            map[string]any
	lock                 *sync.Mutex
	mqttClient           *mqttclient.MqttClient
	plainMetrics         *metrics.PlainMetrics
	statusUpdateInterval int
}

func NewEngine(mqttClient *mqttclient.MqttClient, pm *metrics.PlainMetrics, statusUpdateInterval int) *Engine {
	return &Engine{map[string]any{}, &sync.Mutex{}, mqttClient, pm, statusUpdateInterval}
}

func (e *Engine) Subscribe(mqttListenTopics []string) {
	for _, topic := range mqttListenTopics {
		if err := e.mqttClient.Subscribe(topic, e.messageProcessor); err != nil {
			log.Fatalf("can't subsctibe to: %s", topic)
		}
	}
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
			log.Printf("scheduling status updates for: %s", source)
			e.scheduled[source] = true
			ticker := time.NewTicker(time.Duration(e.statusUpdateInterval) * time.Second)
			go func() {
				for {
					select {
					case <-ticker.C:
						log.Printf("sending status update request command: %s", target)
						if err := e.mqttClient.SendCommand(target, ""); err != nil {
							log.Fatalf("can't send message to: %s", target)
						}
					}
				}
			}()
		}
	}
}
