// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import "google.golang.org/grpc"

func (s *server) registerGRPC(_ *grpc.Server) {
	// Auth gRPC handlers register here when internal/auth is implemented.
}
