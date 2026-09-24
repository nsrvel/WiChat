// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package server

import (
	"encoding/json"
	"net/http"

	"github.com/wichat/wichat/backend/wi-shared/infra/logger"
	"github.com/wichat/wichat/backend/wi-shared/microservices"
	responsehttp "github.com/wichat/wichat/backend/wi-shared/response/http"
)

func (s *server) devAuthPing(w http.ResponseWriter, r *http.Request) {
	ok, err := s.authClient.Ping(r.Context())
	if err != nil {
		responsehttp.WriteError(w, r, err, responsehttp.WithLogger(logger.HTTPWriteErrorLogger(s.log)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":         ok,
		"request_id": microservices.RequestIDFromContext(r.Context()),
	})
}
