package monitor

import (
	"int-status/internal"
	"time"
)

// ServiceStatusChecker defines the interface for checking service status.
type ServiceStatusChecker interface {
	GetTargetServiceConf() internal.ServiceConf
	CheckStatus(timeout time.Duration) internal.Status
}
