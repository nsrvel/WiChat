// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package config

import (
	sharedconfig "github.com/wichat/wichat/backend/wi-shared/config"
	"github.com/wichat/wichat/backend/wi-shared/microservices"
)

// Config is api-gateway service configuration.
type Config struct {
	Env           string                 `mapstructure:"env"`
	LogLevel      string                 `mapstructure:"log_level"`
	Port          string                 `mapstructure:"port"`
	AuthHost      string                 `mapstructure:"ms_auth_host"`
	UserHost      string                 `mapstructure:"ms_user_host"`
	Microservices microservices.Settings `mapstructure:",squash"`
}

// Defaults are applied when env / .env do not set a value.
var Defaults = Config{
	Env:           "development",
	LogLevel:      "info",
	Port:          "3000",
	AuthHost:      "localhost:3001",
	UserHost:      "localhost:3002",
	Microservices: microservices.DefaultSettings,
}

// Load reads configuration from defaults, optional envFile, and environment.
func Load(envFile string) (Config, error) {
	var cfg Config

	if err := sharedconfig.Load(envFile, &cfg, Defaults); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// HTTPListenAddr returns the address for net/http.Server (e.g. ":3000").
func (c Config) HTTPListenAddr() string {
	return sharedconfig.ListenAddress(c.Port)
}
