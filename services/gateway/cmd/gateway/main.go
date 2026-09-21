package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nsrvel/WiChat/services/gateway/internal/httpx"
)

func main() {
	addr := envOr("WICHAT_GATEWAY_ADDR", ":8080")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpx.Health("gateway"))
	mux.HandleFunc("GET /api/v1/health", httpx.Health("gateway"))

	log.Printf("gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
