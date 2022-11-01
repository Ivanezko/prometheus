package telemetry

import (
	"net"
	"strings"
)

const (
	// metricsPrefix is the sub-domain on which the metrics api will be served.
	metricsPrefix      = "metrics."
	defaultMetricsAddr = "metrics.cyolo.io"
	dot                = '.'
)

// Addr returns the metrics api address.
func Addr(addr string) string {
	// addr is an ip address, we return the default metrics api address:
	if ip := net.ParseIP(addr); ip != nil {
		return defaultMetricsAddr
	}
	// if the address is already an metrics api address, return it:
	if strings.HasPrefix(addr, metricsPrefix) {
		return addr
	}
	// addr is not an metrics api address, strip the first sub-domain off of addr:
	if dot := strings.IndexByte(addr, '.'); dot >= 0 {
		addr = addr[dot+1:]
	}
	// concatenate the address with the metrics prefix:
	return metricsPrefix + addr
}

// tenantSuffix returns the domain suffix for the tenantSuffix including port.
func tenantSuffix(sni string) (res string) {
	if firstDot := strings.IndexByte(sni, dot); firstDot > 0 {
		res = sni[firstDot+1:]
	}

	return
}

func TenantSuffix(sni string) (res string) {
	return tenantSuffix(sni)
}
