// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DialOptions configures a shared outbound gRPC connection.
type DialOptions struct {
	// Insecure uses plaintext gRPC (default for local dev).
	Insecure bool
}

// DefaultDialOptions enables insecure transport for development compose setups.
var DefaultDialOptions = DialOptions{Insecure: true}

// Dial opens a client connection with request ID propagation. Pass otel.GRPCClientDialOptions() via extra.
func Dial(target string, opts DialOptions, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	if target == "" {
		return nil, fmt.Errorf("microservices/grpc: dial target is empty")
	}

	dialOpts := make([]grpc.DialOption, 0, len(extra)+2)
	dialOpts = append(dialOpts, extra...)
	dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(UnaryClientRequestID()))
	if opts.Insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("microservices/grpc: dial %s: %w", target, err)
	}

	return conn, nil
}
