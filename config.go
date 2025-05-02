package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General   GeneralConfig `yaml:"general"`
	Endpoint1 EndpointSpec  `yaml:"endpoint1"`
	Endpoint2 EndpointSpec  `yaml:"endpoint2"`
}

type GeneralConfig struct {
	Timeout int               `yaml:"timeout"`
	Headers map[string]string `yaml:"headers"`
}

type EndpointSpec struct {
	URL      string            `yaml:"url"`
	Headers  map[string]string `yaml:"headers"`
	JSONPath string            `yaml:"jsonpath"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	// Set default values
	if config.General.Timeout <= 0 {
		config.General.Timeout = 10 // Default timeout 10 seconds
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Endpoint1.URL == "" {
		return fmt.Errorf("endpoint1 URL is required")
	}

	if config.Endpoint2.URL == "" {
		return fmt.Errorf("endpoint2 URL is required")
	}

	if config.Endpoint1.JSONPath == "" {
		return fmt.Errorf("endpoint1 JSONPath is required")
	}

	if config.Endpoint2.JSONPath == "" {
		return fmt.Errorf("endpoint2 JSONPath is required")
	}

	return nil
}

// Returns the configured timeout duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.General.Timeout) * time.Second
}
