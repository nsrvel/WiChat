// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package config

// Config is ms-auth service configuration.
type Config struct {
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
	GRPCAddr string `mapstructure:"grpc_addr"`
}

// Defaults are applied when env / .env do not set a value.
var Defaults = Config{
	Env:      "development",
	LogLevel: "info",
	GRPCAddr: ":50051",
}
