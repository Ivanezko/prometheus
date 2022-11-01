package telemetry

import (
	ioPrometheusClient "github.com/prometheus/client_model/go"
	"sync"
)

type Vector struct {
	mux        *sync.Mutex
	labelsHash uint64
	Labels     []string
	Value      float64
	Values     []float64
	Observe    func() float64 `json:"-"`
}

func newVector(labelsHash uint64, labels ...string) *Vector {
	return &Vector{
		mux:        &sync.Mutex{},
		labelsHash: labelsHash,
		Labels:     labels,
	}
}

// Drain returns metric copy and resets the original
func (v *Vector) Drain() *Vector {
	v.mux.Lock()
	defer v.mux.Unlock()
	vv := newVector(v.labelsHash, v.Labels...)
	vv.Value = v.Value
	v.Value = 0
	if len(v.Values) != 0 {
		vv.Values = make([]float64, len(v.Values))
		copy(vv.Values, v.Values)
		v.Values = nil
	}
	return vv
}

func (v *Vector) collect(m *Metric, val float64) {
	if val < 0 {
		return
	}
	v.mux.Lock()
	defer v.mux.Unlock()
	switch m.Type {
	case ioPrometheusClient.MetricType_COUNTER:
		v.Value += val
	case ioPrometheusClient.MetricType_HISTOGRAM, ioPrometheusClient.MetricType_SUMMARY:
		v.Values = append(v.Values, val)
	default:
		v.Value = val
	}
}
