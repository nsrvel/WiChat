// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"encoding/json"
	"net/http"

	"github.com/wichat/wichat/backend/gateway/internal/middleware"
	"github.com/wichat/wichat/backend/wi-shared/metrics"
	"github.com/wichat/wichat/backend/wi-shared/ops"
)

func (s *server) buildHandler() (http.Handler, error) {
	mux := http.NewServeMux()
	s.mapHandlers(mux)

	root := middleware.Chain(mux)

	httpMetrics, err := metrics.NewHTTPMetrics(metrics.HTTPOptions{
		Service:  "gateway",
		Registry: s.reg,
	})
	if err != nil {
		return nil, err
	}
	return httpMetrics.Middleware(root), nil
}

func (s *server) mapHandlers(mux *http.ServeMux) {
	ops.Register(mux, s.reg)

	// Auth routes register here when internal/auth/delivery/http exists.

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"service": "gateway"})
	})
}
