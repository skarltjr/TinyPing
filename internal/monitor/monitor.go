package monitor

import (
	"int-status/internal"
	"net/http"
	"time"
)

// StatusMonitor implements ServiceStatusChecker for HTTP-based services.
type StatusMonitor struct {
	service internal.ServiceConf
}

// NewServiceChecker creates a new StatusMonitor instance.
func NewServiceChecker(service internal.ServiceConf) *StatusMonitor {
	return &StatusMonitor{service: service}
}

// GetTargetServiceConf returns the service configuration for the monitor.
func (h *StatusMonitor) GetTargetServiceConf() internal.ServiceConf {
	return h.service
}

// CheckStatus performs a status monitor for the HTTP service.
func (h *StatusMonitor) CheckStatus(timeout time.Duration) internal.Status {
	client := http.Client{
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Get(h.service.API.URL)
	latency := time.Since(start).Milliseconds()

	status := "UP"
	if err != nil || resp.StatusCode >= 500 {
		status = "DOWN"
	}

	return internal.Status{
		Service:   h.service.Name,
		Timestamp: time.Now(),
		Status:    status,
		Latency:   latency,
	}
}
