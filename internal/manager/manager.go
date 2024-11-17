package manager

import (
	"github.com/sirupsen/logrus"
	"int-status/internal"
	"int-status/internal/monitor"
	"int-status/internal/storage"
	"runtime"
	"sync"
	"time"
)

// ServiceManager manages multiple services and their status checks.
type ServiceManager struct {
	checkers []monitor.ServiceStatusChecker
	storage  storage.Storage
}

// NewServiceManager initializes the ServiceManager with a list of services.
func NewServiceManager(services []internal.ServiceConf, storage storage.Storage) *ServiceManager {
	checkers := make([]monitor.ServiceStatusChecker, len(services))
	for i, service := range services {
		checkers[i] = monitor.NewServiceChecker(service)
	}
	return &ServiceManager{checkers: checkers, storage: storage}
}

// StartMonitoring begins periodic status checks for all services
func (m *ServiceManager) StartMonitoring(interval time.Duration) {
	maxGoroutines := runtime.NumCPU()*2 + 10
	guard := make(chan struct{}, maxGoroutines)

	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(nextMinute))

		statusChannel := make(chan internal.Status)
		var wg sync.WaitGroup

		for _, currentMonitor := range m.checkers {
			wg.Add(1)

			go func(monitor monitor.ServiceStatusChecker) {
				defer wg.Done()

				guard <- struct{}{}
				status := monitor.CheckStatus(3 * time.Second)
				statusChannel <- status
				<-guard
			}(currentMonitor)
		}

		go func() {
			wg.Wait()
			close(statusChannel)
		}()

		var statuses []internal.Status
		for status := range statusChannel {
			statuses = append(statuses, status)
		}

		err := m.storage.UpdateHistory(statuses)
		if err != nil {
			logrus.Errorf("Error saving statuses to storage: %v", err)
		}
	}
}

func (m *ServiceManager) GetDailyServiceStatus() (map[string][]internal.Status, error) {
	servicesStatusMap := make(map[string][]internal.Status)
	for _, checker := range m.checkers {
		name := checker.GetTargetServiceConf().Name
		history, err := m.storage.GetDailyHistory(name)
		if err != nil {
			return nil, err
		}
		servicesStatusMap[name] = history
	}

	return servicesStatusMap, nil
}

func (m *ServiceManager) GetDailyIncidents() (map[string][]internal.Incident, error) {
	incidentsMap := make(map[string][]internal.Incident)

	for _, checker := range m.checkers {
		name := checker.GetTargetServiceConf().Name
		incidents, err := m.storage.GetDailyIncidents(name)
		if err != nil {
			return incidentsMap, err
		}
		if len(incidents) > 0 {
			incidentsMap[name] = incidents
		}
	}

	return incidentsMap, nil
}
