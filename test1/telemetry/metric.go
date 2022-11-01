package telemetry

import (
	"fmt"
	ioPrometheusClient "github.com/prometheus/client_model/go"
	"hash/fnv"
	"sync"
)

const (
	MetricTypeCounter   = ioPrometheusClient.MetricType_COUNTER
	MetricTypeGauge     = ioPrometheusClient.MetricType_GAUGE
	MetricTypeSummary   = ioPrometheusClient.MetricType_SUMMARY
	MetricTypeUntyped   = ioPrometheusClient.MetricType_UNTYPED
	MetricTypeHistogram = ioPrometheusClient.MetricType_HISTOGRAM
)

type Metric struct {
	mux        *sync.Mutex
	Name       string
	Type       ioPrometheusClient.MetricType
	LabelNames []string
	Vectors    map[uint64]*Vector
}

func newMetric(name string, typ ioPrometheusClient.MetricType, labelNames ...string) *Metric {
	return &Metric{
		mux:        &sync.Mutex{},
		Name:       name,
		Type:       typ,
		LabelNames: labelNames,
		Vectors:    make(map[uint64]*Vector),
	}
}

// Drain returns metric copy and resets the original
func (m *Metric) Drain() *Metric {
	m.mux.Lock()
	defer m.mux.Unlock()
	mm := newMetric(m.Name, m.Type, m.LabelNames...)
	for _, v := range m.Vectors {
		if v.Observe != nil {
			v.collect(m, v.Observe())
		}
		mm.Vectors[v.labelsHash] = v.Drain()
	}
	return mm
}

func (m *Metric) getVector(labels ...string) *Vector {
	labelsHash := m.hashStrings(labels...)
	if _, ok := m.Vectors[labelsHash]; !ok {
		m.Vectors[labelsHash] = newVector(labelsHash, labels...)
	}
	return m.Vectors[labelsHash]
}

// Collect saves value to the metrics storage
func (m *Metric) Collect(val float64, labels ...string) {
	if val < 0 {
		return
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.getVector(labels...).collect(m, val)
}

// RegisterObserver registers callback with will be called to get value and save it to storage
func (m *Metric) RegisterObserver(o func() float64, labels ...string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.getVector(labels...).Observe = o
}

func (_ *Metric) hashStrings(labels ...string) uint64 {
	var str string
	for _, label := range labels {
		str = fmt.Sprintf("%s:%s", str, label)
	}
	algorithm := fnv.New64a()
	_, _ = algorithm.Write([]byte(str))
	return algorithm.Sum64()
}
