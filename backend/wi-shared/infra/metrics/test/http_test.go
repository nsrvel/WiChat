// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/wichat/wichat/backend/wi-shared/infra/metrics"
)

func TestRoutePatternNormalizesID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/12345/profile", nil)
	got := metrics.RoutePattern(r)
	if got != "/api/v1/users/{id}/profile" {
		t.Fatalf("got %q", got)
	}
}

func TestHTTPMetricsMiddlewareRecords(t *testing.T) {
	reg := prometheus.NewRegistry()
	m, err := metrics.NewHTTPMetrics(metrics.HTTPOptions{Service: "test", Registry: reg})
	if err != nil {
		t.Fatal(err)
	}

	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/99", nil)
	handler.ServeHTTP(rec, req)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, mf := range mfs {
		if mf.GetName() != "wichat_http_requests_total" {
			continue
		}
		for _, metric := range mf.GetMetric() {
			labels := labelMap(metric)
			if labels["status_class"] == "4xx" && strings.Contains(labels["route"], "{id}") {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("metrics: %+v", mfs)
	}
}

func TestNewHTTPMetricsDuplicateRegistry(t *testing.T) {
	reg := prometheus.NewRegistry()
	opts := metrics.HTTPOptions{Service: "test", Registry: reg}
	if _, err := metrics.NewHTTPMetrics(opts); err != nil {
		t.Fatal(err)
	}
	if _, err := metrics.NewHTTPMetrics(opts); err == nil {
		t.Fatal("expected duplicate register error")
	}
}

func labelMap(m *dto.Metric) map[string]string {
	out := make(map[string]string)
	for _, lp := range m.GetLabel() {
		out[lp.GetName()] = lp.GetValue()
	}
	return out
}
