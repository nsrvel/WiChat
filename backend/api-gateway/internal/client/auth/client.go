// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
	"github.com/wichat/wichat/backend/wi-shared/infra/otel"
	"github.com/wichat/wichat/backend/wi-shared/microservices"
	msgrpc "github.com/wichat/wichat/backend/wi-shared/microservices/grpc"
	authv1 "github.com/wichat/wichat/backend/wi-shared/models/gen/auth/v1"
	"google.golang.org/grpc"
)

const readyOverallTimeout = 2 * time.Second

// Client calls ms-auth over gRPC.
type Client struct {
	api      authv1.AuthServiceClient
	conn     *grpc.ClientConn
	log      logger.Logger
	callCfg  microservices.CallConfig
	readyCfg microservices.CallConfig
}

// New dials ms-auth and returns a client. Caller must Close when shutting down.
func New(target string, log logger.Logger, settings microservices.Settings) (*Client, error) {
	conn, err := msgrpc.Dial(target, msgrpc.DefaultDialOptions, otel.GRPCClientDialOptions()...)
	if err != nil {
		return nil, err
	}
	breaker := microservices.NewBreaker(settings)
	callCfg := microservices.CallConfigFromSettings(settings, microservices.RetryIdempotent)
	callCfg.Breaker = breaker

	return &Client{
		api:     authv1.NewAuthServiceClient(conn),
		conn:    conn,
		log:     log,
		callCfg: callCfg,
		readyCfg: microservices.CallConfig{
			PerAttemptTimeout: readyOverallTimeout,
			OverallTimeout:    readyOverallTimeout,
			MaxRetries:        0,
			RetryPolicy:       microservices.RetryNone,
			Breaker:           breaker,
		},
	}, nil
}

// Close releases the gRPC connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// ReadyCheck probes ms-auth for gateway /ready (short timeout, no retry).
func (c *Client) ReadyCheck(ctx context.Context) error {
	ok, err := c.ping(ctx, c.readyCfg)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("auth: ping: not ok")
	}
	return nil
}

// Ping checks reachability of ms-auth.
func (c *Client) Ping(ctx context.Context) (bool, error) {
	return c.ping(ctx, c.callCfg)
}

func (c *Client) ping(ctx context.Context, cfg microservices.CallConfig) (bool, error) {
	var resp *authv1.PingResponse
	err := microservices.CallUnary(ctx, cfg, c.log, authv1.AuthService_Ping_FullMethodName, func(callCtx context.Context) error {
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
