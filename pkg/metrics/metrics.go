package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type Metrics struct {
	pm                        *PlainMetrics
	removeWhenInactiveMinutes int
	gauges                    map[string]*prometheus.GaugeVec
}

func NewMetrics(pm *PlainMetrics, removeWhenInactiveMinutes int) *Metrics {
	m := &Metrics{pm, removeWhenInactiveMinutes, map[string]*prometheus.GaugeVec{}}
	m.scheduleCleanup()
	return m
}

func (m *Metrics) Refresh() {
	m.pm.lock.Lock()
	defer m.pm.lock.Unlock()
	for source, sourceMetrics := range m.pm.metrics {
		statusTopic, topicExists := sourceMetrics.Metrics["status_topic"].(string)
		statusNetHostname, hostnameExists := sourceMetrics.Metrics["status_net_hostname"].(string)
		statusDeviceName, deviceNameExists := sourceMetrics.Metrics["status_device_name"].(string)
		for pmk, pmv := range sourceMetrics.Metrics {
			if float, ok := pmv.(float64); ok && topicExists && hostnameExists && deviceNameExists {
				if _, ok := m.gauges[pmk]; !ok {
					m.gauges[pmk] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name:      pmk,
							Namespace: "tasmota",
						},
						[]string{"source", "status_topic", "status_net_hostname", "status_device_name"},
					)
					prometheus.MustRegister(m.gauges[pmk])
				}
				m.gauges[pmk].WithLabelValues(source, statusTopic, statusNetHostname, statusDeviceName).Set(float)
			}
		}
	}
}

func (m *Metrics) scheduleCleanup() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.pm.lock.Lock()
				for source, sourceMetrics := range m.pm.metrics {
					if time.Since(sourceMetrics.lastAccessed) > time.Duration(m.removeWhenInactiveMinutes)*time.Minute {
						statusTopic, topicExists := sourceMetrics.Metrics["status_topic"].(string)
						statusNetHostname, hostnameExists := sourceMetrics.Metrics["status_net_hostname"].(string)
						statusDeviceName, deviceNameExists := sourceMetrics.Metrics["status_device_name"].(string)
						for pmk, _ := range sourceMetrics.Metrics {
							if topicExists && hostnameExists && deviceNameExists {
								if gauge, ok := m.gauges[pmk]; ok {
									gauge.DeleteLabelValues(source, statusTopic, statusNetHostname, statusDeviceName)
								}
							}
						}
						delete(m.pm.metrics, source)
						log.Println("removed inactive source: " + source)
					}
				}
				m.pm.lock.Unlock()
			}
		}
	}()
	log.Println("scheduled metrics cleanup")
}
