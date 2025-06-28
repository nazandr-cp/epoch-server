package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		// TODO: Add database configuration
	} `yaml:"database"`

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`

	Ethereum struct {
		// TODO: Add Ethereum client configuration
	} `yaml:"ethereum"`

	Subgraph struct {
		// TODO: Add subgraph endpoint configuration
	} `yaml:"subgraph"`

	Scheduler struct {
		Interval time.Duration `yaml:"interval"`
	} `yaml:"scheduler"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
