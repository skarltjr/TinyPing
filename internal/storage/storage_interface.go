package storage

import "int-status/internal"

type Storage interface {
	GetDailyHistory(service string) ([]internal.Status, error)
	GetDailyIncidents(service string) ([]internal.Incident, error)
	UpdateHistory(statuses []internal.Status) error
}
