// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wichat/wichat/backend/gateway/internal/config"
	"github.com/wichat/wichat/backend/wi-shared/logger"
)

// Server runs the public HTTP API.
type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	cfg config.Config
	log logger.Logger
	reg *prometheus.Registry
}

// NewServer constructs the gateway composition root.
func NewServer(cfg config.Config, log logger.Logger, reg *prometheus.Registry) Server {
	return &server{cfg: cfg, log: log, reg: reg}
}
