// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	authclient "github.com/wichat/wichat/backend/api-gateway/internal/client/auth"
)

const httpShutdownTimeout = 5 * time.Second

func (s *server) Run(ctx context.Context) error {
	// Clients
	authClient, err := authclient.New(s.cfg.AuthGRPCAddr, s.log, s.cfg.Microservices)
	if err != nil {
		return err
	}
	s.authClient = authClient

	// Handler
	handler, err := s.buildHandler()
	if err != nil {
		_ = authClient.Close()
		return err
	}

	// HTTP server
	httpSrv := &http.Server{
		Addr:              s.cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Listen
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("http listening", "addr", s.cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http serve: %w", err)
		}
	}()

	// Shutdown
	shutdown := func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
		defer cancel()
		_ = httpSrv.Shutdown(stopCtx)
		if s.authClient != nil {
			_ = s.authClient.Close()
		}
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
