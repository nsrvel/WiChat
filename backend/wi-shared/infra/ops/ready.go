// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package ops

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// DefaultReadyTimeout bounds how long /ready may spend on dependency checks.
const DefaultReadyTimeout = 2 * time.Second

// ReadyCheck probes whether the process can accept traffic.
type ReadyCheck func(ctx context.Context) error

// ReadyHandler runs checks with DefaultReadyTimeout; 503 when any check fails.
func ReadyHandler(checks ...ReadyCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), DefaultReadyTimeout)
		defer cancel()

		for _, check := range checks {
			if check == nil {
				continue
			}
			if err := check(ctx); err != nil {
				writeNotReady(w)
				return
			}
		}

		Health(w, r)
	}
}

func writeNotReady(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "not_ready"})
}
