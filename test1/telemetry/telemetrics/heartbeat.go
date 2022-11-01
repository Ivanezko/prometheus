package telemetrics

import (
	"main/test1/telemetry"
)

type HeartbeatTelemetry struct {
	idacID string
}

func NewHeartbeatWatcher(s *telemetry.Service) {
	m := s.AddMetric("heartbeat", telemetry.MetricTypeCounter, "myLabelName")
	m.RegisterObserver(func() float64 {
		return 123
	}, "observed")
}

func NewHeartbeatMetric(srv *telemetry.Service) *telemetry.Metric {
	m := srv.AddMetric("heartbeat", telemetry.MetricTypeCounter, "myLabelName")
	// m.Collect(1, "lb2val", "lb3val")
	return m
}
