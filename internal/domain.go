package internal

import "time"

// Status represents the real-time status of a service.
// @field Service   The name of the service.
// @field Timestamp The timestamp when the status was recorded.
// @field Status    The current status of the service (e.g., "UP", "DOWN").
// @field Latency   The response time in milliseconds.
type Status struct {
	Service   string
	Timestamp time.Time
	Status    string
	Latency   int64
}

// Incident represents a period of service downtime.
// @field Service   The name of the service that experienced the incident.
// @field StartTime The time when the service went down.
// @field EndTime   The time when the service recovered.
type Incident struct {
	Service   string
	StartTime time.Time
	EndTime   time.Time
}

// ServiceConf represents a single service configuration.
// @field Name        The name of the service.
// @field Description A brief description of the service.
// @field API         details for the service, including method and URL.
type ServiceConf struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	API         struct {
		Method string `yaml:"method"`
		URL    string `yaml:"url"`
	} `yaml:"api"`
}
