package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func recordMetrics() {
	t := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		for {
			_, ok := <-ticker.C
			if !ok {
				fmt.Println("tick not ok")
				return
			}
			fmt.Printf("%+v\n", ticker)
		}
	}(t)
}

func main() {
	recordMetrics()
	fmt.Println("WWW")

	http.Handle("/metrics", promhttp.Handler())
	// http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.ListenAndServe(":2112", nil)
}
