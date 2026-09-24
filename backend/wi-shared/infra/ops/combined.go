// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package ops

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// CombinedServer serves gRPC and ops HTTP (/health, /ready, /metrics) on one TCP port.
type CombinedServer struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

// NewCombinedServer multiplexes grpcSrv and ops routes on a single listener (cleartext h2c for gRPC).
func NewCombinedServer(grpcSrv *grpc.Server, reg prometheus.Gatherer, readyChecks ...ReadyCheck) *CombinedServer {
	mux := http.NewServeMux()
	RegisterWithReady(mux, reg, readyChecks...)

	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	return &CombinedServer{
		grpcServer: grpcSrv,
		httpServer: &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					grpcSrv.ServeHTTP(w, r)
					return
				}
				mux.ServeHTTP(w, r)
			}),
			Protocols:         protocols,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Serve accepts connections on lis until closed or error.
func (s *CombinedServer) Serve(lis net.Listener) error {
	return s.httpServer.Serve(lis)
}

// Shutdown stops gRPC gracefully then the HTTP server.
func (s *CombinedServer) Shutdown(ctx context.Context) error {
	if s.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			s.grpcServer.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-ctx.Done():
			s.grpcServer.Stop()
		}
	}
	return s.httpServer.Shutdown(ctx)
}
