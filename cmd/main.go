// Package main is the entry point for the tokenbridge-k8s-sidecar.
// This sidecar container periodically exchanges a Kubernetes service account token
// for another access token using a configured token exchange endpoint (tokenbridge),
// and writes the result to a file for use by the main application container.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hupe1980/tokenbridge-k8s-sidecar/internal/config"
	"github.com/hupe1980/tokenbridge-k8s-sidecar/internal/token"
)

func main() {
	// Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create a channel to listen for termination signals (SIGINT, SIGTERM).
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("starting token refresher...")

	// Run the token refresher loop until a termination signal is received.
	if err := token.RunRefresher(cfg, stopCh); err != nil {
		log.Fatalf("token refresher failed: %v", err)
	}
}
