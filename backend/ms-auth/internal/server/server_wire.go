// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	authgrpc "github.com/wichat/wichat/backend/ms-auth/internal/auth/delivery/grpc"
	authv1 "github.com/wichat/wichat/backend/wi-shared/models/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (s *server) registerGRPC(grpcSrv *grpc.Server, healthSrv *health.Server) {
	// Health
	healthpb.RegisterHealthServer(grpcSrv, healthSrv)

	// Auth
	authv1.RegisterAuthServiceServer(grpcSrv, authgrpc.NewServer())
}
