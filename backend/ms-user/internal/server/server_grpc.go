// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
	"github.com/wichat/wichat/backend/wi-shared/infra/otel"
	msgrpc "github.com/wichat/wichat/backend/wi-shared/microservices/grpc"
	"google.golang.org/grpc"
)

func newGRPCServer(log logger.Logger) *grpc.Server {
	return msgrpc.NewServer(log, otel.GRPCServerOptions()...)
}
