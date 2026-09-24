// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wichat/wichat/backend/api-gateway/internal/middleware"
	"github.com/wichat/wichat/backend/wi-shared/infra/metrics"
	"github.com/wichat/wichat/backend/wi-shared/infra/ops"
)

func (s *server) buildHandler() (http.Handler, error) {
	// Routes
	mux := http.NewServeMux()
	s.mapHandlers(mux)

	// Middleware
	root := middleware.Chain(mux)

	// HTTP metrics
	httpMetrics, err := metrics.NewHTTPMetrics(metrics.HTTPOptions{
		Service:  "api-gateway",
		Registry: s.reg,
	})
	if err != nil {
		return nil, err
	}

	return httpMetrics.Middleware(root), nil
}

func (s *server) mapHandlers(mux *http.ServeMux) {
	// Ops
	var readyChecks []ops.ReadyCheck
	if s.authClient != nil {
		readyChecks = append(readyChecks, s.authClient.ReadyCheck)
	}
	ops.RegisterWithReady(mux, s.reg, readyChecks...)

	// APIs
	// Auth routes register here when internal/auth/delivery/http exists.
	if !strings.EqualFold(s.cfg.Env, "production") && s.authClient != nil {
		mux.HandleFunc("GET /api/v1/dev/auth-ping", s.devAuthPing)
	}

	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"service": "api-gateway"})
	})
}
