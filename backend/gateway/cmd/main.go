// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wichat/wichat/backend/gateway/internal/config"
	"github.com/wichat/wichat/backend/gateway/internal/server"
	"github.com/wichat/wichat/backend/wi-shared/logger"
	"github.com/wichat/wichat/backend/wi-shared/metrics"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Printf("config: %v", err)
		os.Exit(1)
	}

	appLog := logger.New(logger.Options{
		Service: "gateway",
		Env:     cfg.Env,
		Level:   cfg.LogLevel,
	})

	reg := metrics.NewRegistry()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appLog.Info("gateway starting", "http_addr", cfg.HTTPAddr)

	srv := server.NewServer(cfg, appLog, reg)

	if err := srv.Run(ctx); err != nil {
		appLog.Error("gateway stopped", "err", err)
		os.Exit(1)
	}
	appLog.Info("gateway shutdown complete")
}
