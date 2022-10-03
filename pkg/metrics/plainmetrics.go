package metrics

import (
	"log"
	"strings"
	"sync"
	"time"
)

type SourceMetrics struct {
	Metrics      map[string]any
	lastAccessed time.Time
}

type PlainMetrics struct {
	metrics map[string]*SourceMetrics
	lock    *sync.Mutex
}

func NewPlainMetrics() *PlainMetrics {
	metrics := map[string]*SourceMetrics{}
	return &PlainMetrics{metrics, &sync.Mutex{}}
}

func (m *PlainMetrics) Update(topic string, data map[string]any) {
	if len(data) > 0 {
		segments := strings.Split(topic, "/")
		source := strings.Join(segments[:len(segments)-1], "/")
		log.Println("received update for: " + source)
		m.lock.Lock()
		defer m.lock.Unlock()

		if _, ok := m.metrics[source]; !ok {
			m.metrics[source] = &SourceMetrics{map[string]any{}, time.Now()}
		}

		topicMetrics := m.metrics[source]
		topicMetrics.lastAccessed = time.Now()
		for k, v := range data {
			topicMetrics.Metrics[k] = v
		}
	}
}
