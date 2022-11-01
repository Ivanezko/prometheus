// Package telemetry-simple
// This service provides container for metrics which can accumulate data and send to router
// 	srv := telemetry.NewService("local.cyolo.io", "idac", idacID, "licXXX", "", "", "")
// basic info: tenant, application(router or idac), exemplarId (idacID for idacs)
// All registered metrics are sending every 10 seconds
// Usage if NOT heavy load:
// srv.RegisterMetric(NewMetric("myMetricName", ioPrometheusClient.MetricType_GAUGE, []string{"myLabel"}, nil)).Add(123)
// usage if heavy load (split register and collect, 6x times faster):
// metric := NewMetric("myMetricName", ioPrometheusClient.MetricType_GAUGE, []string{"myLabel"}, nil)
// srv.RegisterMetric(metric)
// for {metric.Collect(1, myLabelValue)}
// You can provide your function as the last parameter - it will be called just before send, so you can use it to collect gauges-like metics

package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	ioPrometheusClient "github.com/prometheus/client_model/go"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	ActiveUsers        = "active_users"
	SupervisorApproval = "supervisor_approval"
	ActiveSessions     = "sessions"
	ConnectionsTracker = "connections_tracker"
	RecordingSessions  = "recording_sessions"
	Applications       = "applications"
	License            = "license"
	Idac               = "idac"
	Location           = "location"
	Host               = "host"
	Heartbeat          = "heartbeat"
)

func NewService(tenant, application, exemplarId string, secret string, upstream string, upstreamSNI string, routerPort string) *Service {
	return &Service{
		mux:           &sync.Mutex{},
		Tenant:        TenantSuffix(tenant),
		Application:   application,
		ExemplarId:    exemplarId,
		url:           fmt.Sprintf("%s:%s", upstreamSNI, routerPort),
		authorization: fmt.Sprintf("Bearer %s", secret),
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
					return net.Dial(network, upstream)
				},
			},
		},
		Metrics: make(map[string]*Metric),
	}
}

type Service struct {
	mux           *sync.Mutex
	Tenant        string
	Application   string
	ExemplarId    string
	url           string
	authorization string
	client        *http.Client
	ticker        *time.Ticker
	Metrics       map[string]*Metric
}

func (s *Service) Start() {
	s.ticker = time.NewTicker(time.Second * 5)
	go func() {
		for {
			_, ok := <-s.ticker.C
			if !ok {
				return
			}
			s.post()
		}
	}()
}

// drain make a copy of main struct and reset the values
func (s *Service) drain() *Service {
	s.mux.Lock()
	defer s.mux.Unlock()

	ss := &Service{
		Tenant:      s.Tenant,
		Application: s.Application,
		ExemplarId:  s.ExemplarId,
		Metrics:     make(map[string]*Metric),
	}
	for _, m := range s.Metrics {
		mm := m.Drain()
		if len(mm.Vectors) > 0 {
			ss.Metrics[mm.Name] = mm
		}
	}
	return ss
}

// sendSamples collect metrics and post to router
func (s *Service) post() {
	ss := s.drain()
	if err := s.Post(ss); err != nil {
		fmt.Println(err)
	}
}

// Post marshal metric payload and post metrics to router
func (s *Service) Post(ss *Service) error {
	body, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	fmt.Printf("sending: %s\n", body)
	return nil
	// return metrics.PostMetricsToRouter(s.client, s.url, s.authorization, bytes.NewBuffer(body))
}

func (s *Service) Stop() {
	s.ticker.Stop()
}

func (s *Service) AddMetric(name string, typ ioPrometheusClient.MetricType, labelNames ...string) *Metric {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.Metrics[name]; ok {
		return s.Metrics[name]
	}
	s.Metrics[name] = newMetric(name, typ, labelNames...)
	return s.Metrics[name]
}
