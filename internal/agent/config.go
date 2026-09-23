package agent

import (
    "os"
    "gopkg.in/yaml.v3"
)

type AgentConfig struct {
    APIKey   string `yaml:"api_key"`
    WorkDir  string `yaml:"work_dir"`
    Token    string `yaml:"token"`
    Projects []int64 `yaml:"projects"`
}

func LoadConfig(path string) (*AgentConfig, error) {
    raw, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg AgentConfig
    err = yaml.Unmarshal(raw, &cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}
