package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type metricsCache struct {
	mu             sync.Mutex
	metricFamilies []*dto.MetricFamily
}

func (m *metricsCache) Post(reg *prometheus.Registry, r io.Reader) error {
	var inFamilies []*dto.MetricFamily
	if err := json.NewDecoder(r).Decode(&inFamilies); err != nil {
		fmt.Errorf(err.Error())
	}
	m.append(inFamilies)
	return nil
}

func (m *metricsCache) Handler(w http.ResponseWriter, r *http.Request) {
	contentType := expfmt.Negotiate(r.Header)
	w.Header().Set("Content-Type", string(contentType))
	enc := expfmt.NewEncoder(w, contentType)
	m.export(enc, w)
	m.clear()
}

func (m *metricsCache) export(enc expfmt.Encoder, w http.ResponseWriter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, mf := range m.metricFamilies {
		if err := enc.Encode(mf); err != nil {
			http.Error(w, "An error has occurred during metrics encoding:\n\n"+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (m *metricsCache) append(newMetricFamilies []*dto.MetricFamily) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metricFamilies = append(m.metricFamilies, newMetricFamilies...)
}

func (m *metricsCache) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metricFamilies = nil
}

func NewMetricCache() *metricsCache {
	return &metricsCache{}
}

func main() {
	flag.Parse()

	fmt.Println("router started")

	reg := prometheus.NewRegistry()
	reg.MustRegister(
	// collectors.NewBuildInfoCollector(),
	// collectors.NewGoCollector(),
	// collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	mCache := NewMetricCache()

	http.HandleFunc("/metrics", mCache.Handler)
	/*	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			// EnableOpenMetrics: true,
			// Pass custom reg
			Registry: reg,
		},
	))*/
	http.HandleFunc("/metrics/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST!")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := mCache.Post(reg, r.Body); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = r.Body.Close()
	})
	fmt.Println("listen on http://127.0.0.1:8030/metrics")
	log.Fatal(http.ListenAndServe(":8030", nil))
}
