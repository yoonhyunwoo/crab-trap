package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Worker  WorkerConfig  `yaml:"worker"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Port    int    `yaml:"port"`
	LogDir  string `yaml:"log_dir"`
}

type WorkerConfig struct {
	MoltbookAPIKey string        `yaml:"moltbook_api_key"`
	ServerURL      string        `yaml:"server_url"`
	Submolt        string        `yaml:"submolt"`
	Interval       time.Duration `yaml:"interval_minutes"`
	OSDetection    bool          `yaml:"os_detection"`
}

type LoggingConfig struct {
	Level        string `yaml:"level"`
	SaveRequests bool   `yaml:"save_requests"`
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

	cfg.Worker.Interval = cfg.Worker.Interval * time.Minute

	return &cfg, nil
}

func LoadDefault() (*Config, error) {
	return Load("config.yaml")
}
