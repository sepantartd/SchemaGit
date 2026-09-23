package agent

import (
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    CloudURL        string `yaml:"cloud_url"`
    APIKey          string `yaml:"api_key"`
    ProjectID       int64  `yaml:"project_id"`
    IntervalSeconds int    `yaml:"interval_seconds"`
    LogLevel        string `yaml:"log_level"`
}

func LoadConfig(path string) (Config, error) {
    var cfg Config
    raw, err := os.ReadFile(path)
    if err != nil {
        return cfg, err
    }
    err = yaml.Unmarshal(raw, &cfg)
    return cfg, err
}
