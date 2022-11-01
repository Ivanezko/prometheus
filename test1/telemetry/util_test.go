package telemetry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddr(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"tcp.local.cyolo.io", "metrics.local.cyolo.io"},
		{"tcp.local.cyolo.io:9090", "metrics.local.cyolo.io:9090"},
		{".local.cyolo.io", "metrics.local.cyolo.io"},
		{"192.56.11.10", defaultMetricsAddr},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, Addr(tt.input))
	}
}

func Test_tenantSuffix(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"chuku.local.cyolo.io", "local.cyolo.io"},
		{"chuku.local.cyolo.io:9090", "local.cyolo.io:9090"},
		{"chuku.buku.local.cyolo.io", "buku.local.cyolo.io"},
		{"tcp.cyolo.io", "cyolo.io"},
		{"something that is not a domain is not valid", ""},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, tenantSuffix(tt.input))
	}
}
