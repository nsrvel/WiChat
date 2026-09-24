// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/wichat/wichat/backend/wi-shared/exception"
	responsehttp "github.com/wichat/wichat/backend/wi-shared/response/http"
)

const requestIDHeader = "X-Request-ID"

// Chain wraps the handler with gateway middleware (order: request ID, recover).
func Chain(next http.Handler) http.Handler {
	return requestID(recoverPanic(next))
}

// Request ID
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}

		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])

	return hex.EncodeToString(b[:])
}

// Recover panic
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				responsehttp.WriteError(w, r, exception.Internal(nil))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
