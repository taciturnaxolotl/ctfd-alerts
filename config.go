package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type CTFdConfig struct {
	ApiBase string `toml:"api_base"`
	ApiKey  string `toml:"api_key"`
}

type NtfyConfig struct {
	ApiBase     string `toml:"api_base"`
	AccessToken string `toml:"acess_token"`
	Topic       string `toml:"topic"`
}

type Config struct {
	Debug           bool       `toml:"debug"`
	User            string     `toml:"user"`
	CTFdConfig      CTFdConfig `toml:"ctfd"`
	NtfyConfig      NtfyConfig `toml:"ntfy"`
	MonitorInterval int        `toml:"interval"`
}

var config *Config

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.CTFdConfig.ApiBase == "" {
		return nil, errors.New("ctfd api_base URL cannot be empty")
	}

	if cfg.CTFdConfig.ApiKey == "" {
		return nil, errors.New("ctfd api_key cannot be empty")
	}

	// Check API key format (should start with ctfd_ followed by 64 hex characters)
	if len(cfg.CTFdConfig.ApiKey) != 69 || !strings.HasPrefix(cfg.CTFdConfig.ApiKey, "ctfd_") {
		return nil, errors.New("ctfd api_key must be in the format ctfd_<64 hex characters> not " + cfg.CTFdConfig.ApiKey)
	}

	if cfg.NtfyConfig.ApiBase == "" {
		return nil, errors.New("ntfy api_base URL cannot be empty")
	}

	if cfg.NtfyConfig.Topic == "" {
		return nil, errors.New("ntfy topic cannot be empty")
	}

	if cfg.User == "" {
		return nil, errors.New("user cannot be empty")
	}

	if cfg.MonitorInterval == 0 {
		cfg.MonitorInterval = 300
		fmt.Println("you haven't set a monitor interval; setting to 300")
	}

	return &cfg, nil
}
