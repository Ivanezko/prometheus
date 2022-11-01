package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"math/rand"
	"net/http"
	"time"
)

var (
	reg     = &prometheus.Registry{}
	counter = &prometheus.CounterVec{}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("start sending metrics")

	reg = prometheus.NewRegistry()
	reg.MustRegister(
	// collectors.NewBuildInfoCollector(),
	// collectors.NewGoCollector(),
	// collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_counter",
		Help: "counter example",
	}, []string{"node", "mylabel"})
	reg.MustRegister(counter)

	counter.WithLabelValues("1", "myval1").Add(1)

	mClient := NewMetricsClient(reg)
	mClient.PostMetricsToRouter()
	counter.WithLabelValues("1", "myval1").Add(4)
	mClient.PostMetricsToRouter()
}

type metricsClient struct {
	reg *prometheus.Registry
}

func NewMetricsClient(reg *prometheus.Registry) *metricsClient {
	return &metricsClient{reg: reg}
}

func (m *metricsClient) PostMetricsToRouter() {
	metricFamilies, _ := reg.Gather()
	metricFamiliesJSON, _ := json.Marshal(metricFamilies)
	resp, _ := http.Post("http://127.0.0.1:8030/metrics/", "application/json", bytes.NewBuffer(metricFamiliesJSON))
	resp.Body.Close()
}
