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

	"github.com/wichat/wichat/backend/wi-shared/ops"
	"google.golang.org/grpc"
)

const grpcShutdownTimeout = 5 * time.Second

func (s *server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("grpc listen %s: %w", s.cfg.GRPCAddr, err)
	}

	grpcSrv := newGRPCServer(s.log)
	s.registerGRPC(grpcSrv)
	opsHTTP := ops.NewHTTPServer(s.cfg.MetricsAddr, s.reg)

	errCh := make(chan error, 2)

	go func() {
		s.log.Info("grpc listening", "addr", s.cfg.GRPCAddr)
		if err := grpcSrv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	go func() {
		s.log.Info("ops http listening", "addr", s.cfg.MetricsAddr)
		if err := opsHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("ops http serve: %w", err)
		}
	}()

	shutdown := func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), grpcShutdownTimeout)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			grpcSrv.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-stopCtx.Done():
			grpcSrv.Stop()
		}

		opsCtx, opsCancel := context.WithTimeout(context.Background(), ops.ShutdownTimeout)
		defer opsCancel()
		_ = opsHTTP.Shutdown(opsCtx)
	}

	select {
	case <-ctx.Done():
		shutdown()
		return nil
	case err := <-errCh:
		shutdown()
		return err
	}
}
