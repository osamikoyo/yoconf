package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HttpPort int    `yaml:"http_port"`
	GrpcPort int    `yaml:"grpc_port"`
	Addr     string `yaml:"addr"`
	RedisUrl string `yaml:"redis_url"`
	CfgDir   string `yaml:"cfg_dir"`
}

func NewConfig(addr string) (*Config, error) {
	cfg := Config{}

	file, err := os.Open(addr)
	if err != nil {
		return nil, fmt.Errorf("failed open config file: %v", err)
	}

	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed decode config: %v", err)
	}

	return &cfg, nil
}
