package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/push"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			for i := 0; i < 100; i++ {
				rnd := float64(i)
				counter.WithLabelValues("1").Add(1)
				gauge.WithLabelValues("1").Set(rnd)
				hist.WithLabelValues("1").Observe(rnd)
				summary.WithLabelValues("1").Observe(rnd)
			}
			// time.Sleep(10 * time.Second)
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(counter).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			counter.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(gauge).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			gauge.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(hist).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			hist.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(summary).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push summary time to Pushgateway:", err)
			}
			summary.Reset()
			time.Sleep(10 * time.Second)
		}

	}()
	/*	go func() {
		for i := 0; ; i++ {
			rnd := float64(i)
			counter.WithLabelValues("2").Add(1)
			gauge.WithLabelValues("2").Set(rnd)

			// hist.Reset()
			hist.WithLabelValues("2").Observe(rnd)
			time.Sleep(10 * time.Second)
		}
	}()*/
	go func() {
		return
		for {
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(counter).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			counter.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(gauge).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			gauge.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(hist).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
			hist.Reset()
			if err := push.New("http://127.0.0.1:8070", "db_backup").
				Collector(summary).
				Grouping("db", "customers").
				Format(expfmt.FmtText).
				Push(); err != nil {
				fmt.Println("Could not push summary time to Pushgateway:", err)
			}
			summary.Reset()
			time.Sleep(5 * time.Second)
		}
	}()
}

var (
	reg     = &prometheus.Registry{}
	counter = &prometheus.CounterVec{}
	gauge   = &prometheus.GaugeVec{}
	hist    = &prometheus.HistogramVec{}
	summary = &prometheus.SummaryVec{}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("start sending metrics")

	reg = prometheus.NewRegistry()
	reg.MustRegister(
		// collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		// collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_counter",
		Help: "counter example",
	}, []string{"node", "mylabel"})
	reg.MustRegister(counter)
	gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "myapp_gauge",
		Help: "gauge example",
	}, []string{"node"})
	reg.MustRegister(gauge)
	/*	hist := promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "myapp_hist",
			Help:    "hist example",
			Buckets: []float64{10, 100, 500},
		}, []string{"node"})
		summary := promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "myapp_summary",
			Help: "hist example",
		}, []string{"node"})*/

	counter.WithLabelValues("1", "myval1").Add(1)
	counter.WithLabelValues("2", "myval2").Add(2)
	gauge.WithLabelValues("1").Set(123)

	metricFamilies, err := reg.Gather()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	// fmt.Printf("%+v", metrics)

	// temp reg to collect converted metrics
	regTmp := prometheus.NewRegistry()
	for _, metricFamily := range metricFamilies {
		fmt.Printf("metricFamily: %+v\n", metricFamily)

		labels := []string{"service"}
		metrics := metricFamily.GetMetric()
		for _, labelPair := range metrics[0].GetLabel() {
			labels = append(labels, *labelPair.Name)
		}

		switch *metricFamily.Type {
		case io_prometheus_client.MetricType_COUNTER:
			newMetricFamily := promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "transfer_" + *metricFamily.Name,
				Help: *metricFamily.Help,
			}, labels)
			regTmp.MustRegister(newMetricFamily)

			for _, metric := range metrics {
				labelValues := []string{"idac"}
				for _, labelPair := range metric.GetLabel() {
					labelValues = append(labelValues, labelPair.GetValue())
				}
				newMetricFamily.WithLabelValues(labelValues...).Add(metric.Counter.GetValue())
			}
		case io_prometheus_client.MetricType_GAUGE:
			newMetricFamily := promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "transfer_" + *metricFamily.Name,
				Help: *metricFamily.Help,
			}, labels)
			regTmp.MustRegister(newMetricFamily)

			for _, metric := range metrics {
				labelValues := []string{"idac"}
				for _, labelPair := range metric.GetLabel() {
					labelValues = append(labelValues, labelPair.GetValue())
				}
				newMetricFamily.WithLabelValues(labelValues...).Set(metric.Gauge.GetValue())
			}
		}

	}

	if err := push.New("http://127.0.0.1:8070", "transfer").
		Gatherer(regTmp).
		Grouping("db", "customers").
		Format(expfmt.FmtText).
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}

	// recordMetrics()

	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			// EnableOpenMetrics: true,
			// Pass custom reg
			Registry: reg,
		},
	))

	http.ListenAndServe(":2112", nil)
}
