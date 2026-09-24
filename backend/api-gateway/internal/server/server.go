// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	authclient "github.com/wichat/wichat/backend/api-gateway/internal/client/auth"
	"github.com/wichat/wichat/backend/api-gateway/internal/config"
	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
)

// Server runs the public HTTP API.
type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	cfg        config.Config
	log        logger.Logger
	reg        *prometheus.Registry
	authClient *authclient.Client
}

// NewServer constructs the api-gateway composition root.
func NewServer(cfg config.Config, log logger.Logger, reg *prometheus.Registry) Server {
	return &server{cfg: cfg, log: log, reg: reg}
}
