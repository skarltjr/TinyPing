package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"int-status/internal"
	"os"
)

// LoadServices loads the YAML configuration file.
func LoadServices(path string) ([]internal.ServiceConf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var services []internal.ServiceConf
	if err := yaml.Unmarshal(data, &services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return services, nil
}
