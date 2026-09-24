// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package auth

import (
	"context"
	"fmt"

	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
	"github.com/wichat/wichat/backend/wi-shared/infra/otel"
	"github.com/wichat/wichat/backend/wi-shared/microservices"
	msgrpc "github.com/wichat/wichat/backend/wi-shared/microservices/grpc"
	authv1 "github.com/wichat/wichat/backend/wi-shared/models/gen/auth/v1"
	"google.golang.org/grpc"
)

// Client calls ms-auth over gRPC.
type Client struct {
	api     authv1.AuthServiceClient
	conn    *grpc.ClientConn
	log     logger.Logger
	callCfg microservices.CallConfig
}

// New dials ms-auth and returns a client. Caller must Close when shutting down.
func New(target string, log logger.Logger, settings microservices.Settings) (*Client, error) {
	conn, err := msgrpc.Dial(target, msgrpc.DefaultDialOptions, otel.GRPCClientDialOptions()...)
	if err != nil {
		return nil, err
	}
	return &Client{
		api:     authv1.NewAuthServiceClient(conn),
		conn:    conn,
		log:     log,
		callCfg: microservices.CallConfigFromSettings(settings, microservices.RetryIdempotent),
	}, nil
}

// Close releases the gRPC connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Ping checks reachability of ms-auth.
func (c *Client) Ping(ctx context.Context) (bool, error) {
	var resp *authv1.PingResponse
	err := microservices.CallUnary(ctx, c.callCfg, c.log, authv1.AuthService_Ping_FullMethodName, func(callCtx context.Context) error {
		out, callErr := c.api.Ping(callCtx, &authv1.PingRequest{})
		resp = out
		return callErr
	})
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("auth: ping: empty response")
	}
	return resp.Ok, nil
}
