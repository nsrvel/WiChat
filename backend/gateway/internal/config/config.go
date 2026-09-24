// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package config

import (
	sharedconfig "github.com/wichat/wichat/backend/wi-shared/config"
	"github.com/wichat/wichat/backend/wi-shared/microservices"
)

// Config is gateway service configuration.
type Config struct {
	Env           string                 `mapstructure:"env"`
	LogLevel      string                 `mapstructure:"log_level"`
	HTTPAddr      string                 `mapstructure:"http_addr"`
	AuthGRPCAddr  string                 `mapstructure:"auth_grpc_addr"`
	Microservices microservices.Settings `mapstructure:",squash"`
}

// Defaults are applied when env / .env do not set a value.
var Defaults = Config{
	Env:           "development",
	LogLevel:      "info",
	HTTPAddr:      ":8080",
	AuthGRPCAddr:  "localhost:50051",
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
