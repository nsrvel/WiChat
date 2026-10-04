// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package ops

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ServiceReadyURL builds an HTTP /ready URL from a host (e.g. localhost:3001 or http://localhost:3001).
func ServiceReadyURL(host string) string {
	h := strings.TrimSpace(host)
	h = strings.TrimSuffix(h, "/")
	if strings.HasPrefix(h, "http://") || strings.HasPrefix(h, "https://") {
		return h + "/ready"
	}
	return "http://" + h + "/ready"
}

// HTTPGetReadyCheck returns a ReadyCheck that GETs url and requires status 200.
func HTTPGetReadyCheck(client *http.Client, url string) ReadyCheck {
	if client == nil {
		client = http.DefaultClient
	}
	return func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("GET %s: %s", url, resp.Status)
		}
		return nil
	}
}
