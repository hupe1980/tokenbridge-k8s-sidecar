package token

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/hupe1980/tokenbridge-k8s-sidecar/internal/config"
)

// RunRefresher continuously refreshes the access token.
// It uses a retry with exponential backoff on failure and gracefully handles shutdown signals.
func RunRefresher(cfg *config.Config, stopCh <-chan os.Signal) error {
	for {
		expiry := refreshOnceWithRetry(cfg, stopCh)

		// Determine next refresh interval
		sleep := cfg.RefreshInterval
		if expiry > 0 && expiry < int64(cfg.RefreshInterval.Seconds()) {
			sleep = time.Duration(expiry) * time.Second
		}

		select {
		case <-time.After(sleep):
		case <-stopCh:
			log.Println("received termination signal during sleep, shutting down...")
			return nil
		}
	}
}

// refreshOnceWithRetry attempts to refresh the token with retry and exponential backoff.
func refreshOnceWithRetry(cfg *config.Config, stopCh <-chan os.Signal) int64 {
	var (
		maxRetries   = 5
		backoff      = 2 * time.Second
		maxBackoff   = 30 * time.Second
		retryAttempt = 0
	)

	for {
		expiresIn, err := refreshOnce(cfg)
		if err == nil {
			log.Println("token refreshed successfully")
			return expiresIn
		}

		log.Printf("token refresh failed (attempt %d): %v", retryAttempt+1, err)

		retryAttempt++

		if retryAttempt >= maxRetries {
			log.Printf("maximum retries (%d) reached, continuing with next cycle", maxRetries)
			return 0
		}

		select {
		case <-time.After(backoff):
			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
		case <-stopCh:
			log.Println("received termination signal during retry, shutting down...")
			return 0
		}
	}
}

// refreshOnce performs a single token exchange and writes the token to file.
func refreshOnce(cfg *config.Config) (int64, error) {
	idToken, err := os.ReadFile(cfg.SATokenPath)
	if err != nil {
		return 0, err
	}

	resp, err := ExchangeToken(cfg.ExchangeURL, string(idToken), cfg.Audience)
	if err != nil {
		return 0, err
	}

	if err := os.WriteFile(cfg.OutputTokenPath, []byte(resp.AccessToken), 0600); err != nil {
		return 0, err
	}

	return resp.ExpiresIn, nil
}
