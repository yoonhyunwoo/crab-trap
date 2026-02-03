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
	IntervalMinutes int          `yaml:"interval_minutes"`
	Interval       time.Duration `yaml:"-"`
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

	cfg.Worker.Interval = time.Duration(cfg.Worker.IntervalMinutes) * time.Minute

	return &cfg, nil
}

func LoadDefault() (*Config, error) {
	return Load("config.yaml")
}
