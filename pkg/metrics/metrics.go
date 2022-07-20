package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type Metrics struct {
	pm     *PlainMetrics
	gauges map[string]*prometheus.GaugeVec
	lock   *sync.Mutex
}

func NewMetrics(pm *PlainMetrics) *Metrics {
	return &Metrics{pm, map[string]*prometheus.GaugeVec{}, &sync.Mutex{}}
}

func (m *Metrics) Refresh() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for source, sourceMetrics := range m.pm.metrics {
		statusTopic, topicExists := sourceMetrics["status_topic"].(string)
		statusNetHostname, hostnameExists := sourceMetrics["status_net_hostname"].(string)
		statusDeviceName, deviceNameExists := sourceMetrics["status_device_name"].(string)
		for pmk, pmv := range sourceMetrics {
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
