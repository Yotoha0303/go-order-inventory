package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig `yaml:"server"`
	MySQL      MySQLConfig  `yaml:"mysql"`
	Redis      RedisConfig  `yaml:"redis"`
	HttpServer HttpServer   `yaml:"http"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MySQLConfig struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`

	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifeTime"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	PingTimeout     time.Duration `yaml:"pingTimeout"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
	DB   int    `yaml:"db"`
}

type HttpServer struct {
	Server HttpServerConfig `yaml:"server"`
}

type HttpServerConfig struct {
	ReadTimeOut       time.Duration `yaml:"readTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
	MaxHeaderBytesKib int           `yaml:"maxHeaderBytesKib"`
	Timeout           time.Duration `yaml:"timeout"`
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

	if err = cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
