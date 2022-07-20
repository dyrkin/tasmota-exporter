package metrics

import (
	"strings"
	"sync"
)

type PlainMetrics struct {
	metrics map[string]map[string]any
	lock    *sync.Mutex
}

func NewPlainMetrics() *PlainMetrics {
	return &PlainMetrics{map[string]map[string]any{}, &sync.Mutex{}}
}

func (m *PlainMetrics) Update(topic string, data map[string]any) {
	if len(data) > 0 {
		segments := strings.Split(topic, "/")
		source := strings.Join(segments[:len(segments)-1], "/")
		m.lock.Lock()
		defer m.lock.Unlock()

		if _, ok := m.metrics[source]; !ok {
			m.metrics[source] = map[string]any{}
		}

		topicMetrics := m.metrics[source]
		for k, v := range data {
			topicMetrics[k] = v
		}
	}
}
