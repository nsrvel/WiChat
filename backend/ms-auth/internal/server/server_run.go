// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/wichat/wichat/backend/wi-shared/infra/ops"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const serveShutdownTimeout = 5 * time.Second

func (s *server) Run(ctx context.Context) error {
	// gRPC
	lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.cfg.GRPCAddr, err)
	}

	grpcSrv := newGRPCServer(s.log)
	healthSrv := health.NewServer()
	s.registerGRPC(grpcSrv, healthSrv)

	// Ops + gRPC on one port
	readyCheck := func(ctx context.Context) error {
		resp, err := healthSrv.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
		if err != nil {
			return err
		}
		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			return fmt.Errorf("grpc health: %s", resp.GetStatus())
		}
		return nil
	}
	combined := ops.NewCombinedServer(grpcSrv, s.reg, readyCheck)

	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Listen
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("listening", "addr", s.cfg.GRPCAddr, "grpc", true, "ops_http", true)
		if err := combined.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("serve: %w", err)
		}
	}()

	// Shutdown
	shutdown := func() {
		healthSrv.Shutdown()

		stopCtx, cancel := context.WithTimeout(context.Background(), serveShutdownTimeout)
		defer cancel()
		_ = combined.Shutdown(stopCtx)
	}

	// Wait
	select {
	case <-ctx.Done():
		shutdown()
		return nil
	case err := <-errCh:
		shutdown()
		return err
	}
}
