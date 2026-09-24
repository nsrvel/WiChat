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
	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
	"github.com/wichat/wichat/backend/wi-shared/infra/metrics"
	"github.com/wichat/wichat/backend/wi-shared/infra/otel"
)

func main() {
	// Config
	cfg, err := config.Load(".env")
	if err != nil {
		log.Printf("config: %v", err)
		os.Exit(1)
	}

	// Logger
	appLog := logger.New(logger.Options{
		Service: "gateway",
		Env:     cfg.Env,
		Level:   cfg.LogLevel,
	})

	// Metrics
	reg := metrics.NewRegistry()

	// Tracing
	shutdownTracer, err := otel.InstallTracer(context.Background(), "gateway")
	if err != nil {
		appLog.Error("otel init failed", "err", err)
		os.Exit(1)
	}
	defer func() { _ = shutdownTracer(context.Background()) }()

	// Signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appLog.Info("gateway starting", "http_addr", cfg.HTTPAddr)

	// Server
	srv := server.NewServer(cfg, appLog, reg)
	if err := srv.Run(ctx); err != nil {
		appLog.Error("gateway stopped", "err", err)
		os.Exit(1)
	}

	appLog.Info("gateway shutdown complete")
}
