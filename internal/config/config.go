// Package config handles the configuration loading for the tokenbridge-k8s-sidecar.
// It reads required and optional environment variables and constructs a Config object
// used to control the behavior of the sidecar.
package config

import (
	"errors"
	"os"
	"time"
)

// Config holds all configuration settings required by the token refresher.
// These values are typically loaded from environment variables.
type Config struct {
	// SATokenPath is the file path to the projected Kubernetes service account token.
	SATokenPath string

	// OutputTokenPath is the file path where the exchanged token should be written.
	OutputTokenPath string

	// ExchangeURL is the endpoint URL of the tokenbridge service used for token exchange.
	ExchangeURL string

	// RefreshInterval determines how often the token should be refreshed.
	// If the interval is shorter than the token's expiration, it ensures continuous availability.
	RefreshInterval time.Duration

	// Audience is an optional parameter used when exchanging the token;
	// it represents the intended audience for the token.
	Audience string
}

// Load reads configuration from environment variables and returns a Config instance.
// It returns an error if any required environment variables are missing or invalid.
func Load() (*Config, error) {
	saToken := os.Getenv("SA_TOKEN_PATH")
	if saToken == "" {
		return nil, errors.New("SA_TOKEN_PATH is required")
	}

	outputPath := os.Getenv("OUTPUT_TOKEN_PATH")
	if outputPath == "" {
		return nil, errors.New("OUTPUT_TOKEN_PATH is required")
	}

	exchangeURL := os.Getenv("EXCHANGE_URL")
	if exchangeURL == "" {
		return nil, errors.New("EXCHANGE_URL is required")
	}

	intervalStr := os.Getenv("REFRESH_INTERVAL")
	if intervalStr == "" {
		intervalStr = "1h" // Default interval
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	audience := os.Getenv("AUDIENCE") // Optional

	return &Config{
		SATokenPath:     saToken,
		OutputTokenPath: outputPath,
		ExchangeURL:     exchangeURL,
		RefreshInterval: interval,
		Audience:        audience,
	}, nil
}
