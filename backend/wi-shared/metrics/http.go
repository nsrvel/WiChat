// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Garam

package metrics

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HTTPOptions configures RED HTTP metrics.
type HTTPOptions struct {
	Service  string
	Registry *prometheus.Registry
}

// HTTPMetrics records RED metrics; create once per registry via NewHTTPMetrics.
type HTTPMetrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// NewHTTPMetrics registers collectors on opts.Registry (once per registry/service pair).
func NewHTTPMetrics(opts HTTPOptions) (*HTTPMetrics, error) {
	if opts.Registry == nil {
		return nil, errors.New("metrics: registry is required")
	}

	constLabels := prometheus.Labels{}
	if opts.Service != "" {
		constLabels["service"] = opts.Service
	}

	m := &HTTPMetrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "wichat_http_requests_total",
			Help:        "Total HTTP requests",
			ConstLabels: constLabels,
		}, []string{"method", "route", "status_class"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "wichat_http_request_duration_seconds",
			Help:        "HTTP request latency in seconds",
			ConstLabels: constLabels,
			Buckets:     prometheus.DefBuckets,
		}, []string{"method", "route"}),
	}

	if err := opts.Registry.Register(m.requests); err != nil {
		return nil, err
	}
	if err := opts.Registry.Register(m.duration); err != nil {
		return nil, err
	}
	return m, nil
}

// Middleware records RED metrics for requests handled by next.
func (m *HTTPMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		route := RoutePattern(r)
		method := r.Method
		class := statusClass(rw.status)
		m.requests.WithLabelValues(method, route, class).Inc()
		m.duration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
	})
}

// WrapHTTP is a convenience helper that registers metrics and wraps next.
// Prefer NewHTTPMetrics at startup when wrapping multiple handlers on the same registry.
func WrapHTTP(next http.Handler, opts HTTPOptions) (http.Handler, error) {
	m, err := NewHTTPMetrics(opts)
	if err != nil {
		return nil, err
	}
	return m.Middleware(next), nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := r.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("metrics: hijack not supported")
}

func statusClass(code int) string {
	switch {
	case code >= 500:
		return "5xx"
	case code >= 400:
		return "4xx"
	case code >= 300:
		return "3xx"
	default:
		return "2xx"
	}
}

const maxRouteSegments = 12

var (
	uuidLike  = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	numericID = regexp.MustCompile(`^\d+$`)
)

// RoutePattern returns a low-cardinality route label for metrics.
func RoutePattern(r *http.Request) string {
	path := r.URL.Path
	if path == "" {
		return "/"
	}
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) > maxRouteSegments {
		return "unknown"
	}
	for i, seg := range segments {
		if seg == "" {
			continue
		}
		if uuidLike.MatchString(seg) || numericID.MatchString(seg) {
			segments[i] = "{id}"
		}
	}
	return "/" + strings.Join(segments, "/")
}
