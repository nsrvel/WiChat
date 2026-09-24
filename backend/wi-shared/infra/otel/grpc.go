// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package otel

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// GRPCClientDialOptions returns dial options for gRPC client tracing.
func GRPCClientDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

// GRPCServerOptions returns server options for gRPC server tracing.
func GRPCServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
}
