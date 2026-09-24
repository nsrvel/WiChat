// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wichat/wichat/backend/wi-shared/logger"
)

func newGRPCServer(log logger.Logger) *grpc.Server {
	return grpc.NewServer(grpc.ChainUnaryInterceptor(recoveryUnary(log)))
}

func recoveryUnary(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("grpc panic recovered",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		resp, err = handler(ctx, req)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Internal {
				log.Error("grpc internal error", "method", info.FullMethod, "err", err)
			}
		}
		return resp, err
	}
}
