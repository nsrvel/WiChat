// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package ops

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wichat/wichat/backend/wi-shared/metrics"
)

// ShutdownTimeout is the default grace period for ops HTTP shutdown.
const ShutdownTimeout = 5 * time.Second

// Health writes the standard liveness JSON response.
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Register mounts /health and /metrics on mux.
func Register(mux *http.ServeMux, reg prometheus.Gatherer) {
	mux.HandleFunc("/health", Health)
	mux.Handle("/metrics", metrics.Handler(reg))
}

// HTTPServer serves health and Prometheus metrics on a dedicated listen address.
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer builds an ops listener (microservices: metrics_addr).
func NewHTTPServer(addr string, reg prometheus.Gatherer) *HTTPServer {
	// Routes
	mux := http.NewServeMux()
	Register(mux, reg)

	// Server
	return &HTTPServer{
		server: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// ListenAndServe starts the ops HTTP server.
func (s *HTTPServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
