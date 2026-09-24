// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package ops_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wichat/wichat/backend/wi-shared/infra/ops"
)

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	ops.Health(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("body %s", rec.Body.String())
	}
}

func TestReadyOK(t *testing.T) {
	rec := httptest.NewRecorder()
	ops.ReadyHandler(func(ctx context.Context) error { return nil })(
		rec, httptest.NewRequest(http.MethodGet, "/ready", nil),
	)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestReadyFails(t *testing.T) {
	rec := httptest.NewRecorder()
	ops.ReadyHandler(func(ctx context.Context) error { return context.DeadlineExceeded })(
		rec, httptest.NewRequest(http.MethodGet, "/ready", nil),
	)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"not_ready"`) {
		t.Fatalf("body %s", rec.Body.String())
	}
}

func TestRegisterWithReadyRoute(t *testing.T) {
	mux := http.NewServeMux()
	reg := prometheus.NewRegistry()
	ops.RegisterWithReady(mux, reg, func(ctx context.Context) error { return nil })

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("ready status %d", rec.Code)
	}
}

func TestRegisterMetricsRoute(t *testing.T) {
	mux := http.NewServeMux()
	reg := prometheus.NewRegistry()
	ops.Register(mux, reg)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("metrics status %d", rec.Code)
	}
}
