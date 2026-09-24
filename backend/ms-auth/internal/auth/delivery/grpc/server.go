// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package grpc

import (
	"context"

	authv1 "github.com/wichat/wichat/backend/wi-shared/models/gen/auth/v1"
)

// Server implements auth.v1.AuthService.
type Server struct {
	authv1.UnimplementedAuthServiceServer
}

// NewServer constructs the auth gRPC delivery layer.
func NewServer() *Server {
	return &Server{}
}

// Ping reports that the auth service is reachable.
func (s *Server) Ping(_ context.Context, _ *authv1.PingRequest) (*authv1.PingResponse, error) {
	return &authv1.PingResponse{Ok: true}, nil
}
