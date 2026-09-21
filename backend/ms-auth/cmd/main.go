// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package main

import (
	"log"
	"os"

	"github.com/wichat/wichat/backend/ms-auth/internal/config"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Printf("config: %v", err)
		os.Exit(1)
	}

	log.Printf("ms-auth ready env=%s log_level=%s grpc_addr=%s", cfg.Env, cfg.LogLevel, cfg.GRPCAddr)
}
