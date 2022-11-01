package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"main/test1/telemetry"
	"main/test1/telemetry/telemetrics"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	idacID := "XXX1"
	srv := telemetry.NewService("local.cyolo.io", "idac", idacID, "licXXX", "", "", "")
	metric := telemetrics.NewHeartbeatMetric(srv)
	telemetrics.NewHeartbeatWatcher(srv)

	srv.Start()

	go func() {
		for i := 0; ; i++ {
			metric.Collect(1, "fast")
			// time.Sleep(time.Second)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			metric.Collect(1, "slow")
			time.Sleep(time.Millisecond * 100)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":2112", nil)
}
