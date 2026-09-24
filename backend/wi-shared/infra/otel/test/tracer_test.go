// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package otel_test

import (
	"context"
	"testing"

	"github.com/wichat/wichat/backend/wi-shared/infra/otel"
)

func TestInstallTracerNoEndpoint(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	shutdown, err := otel.InstallTracer(context.Background(), "test")
	if err != nil {
		t.Fatalf("InstallTracer: %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
