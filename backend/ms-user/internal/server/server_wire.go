// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (s *server) registerGRPC(grpcSrv *grpc.Server, healthSrv *health.Server) {
	healthpb.RegisterHealthServer(grpcSrv, healthSrv)
}
