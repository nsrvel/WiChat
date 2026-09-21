// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package config

import (
	sharedconfig "github.com/wichat/wichat/backend/wi-shared/config"
)

// Load reads configuration from defaults, optional envFile, and environment.
func Load(envFile string) (Config, error) {
	var cfg Config
	if err := sharedconfig.Load(envFile, &cfg, Defaults); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
