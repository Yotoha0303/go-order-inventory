package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MySQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func LoadEnv() {
	_ = godotenv.Load()
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s failed: %w", path, err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file %s failed: %w", path, err)
	}

	return &cfg, nil
}
