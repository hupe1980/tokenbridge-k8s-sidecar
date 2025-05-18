package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ExchangeResponse represents the structure of the response returned from the
// tokenbridge service after a successful token exchange.
type ExchangeResponse struct {
	// AccessToken is the newly issued token to be used by the consuming application.
	AccessToken string `json:"access_token"`

	// IssuedTokenType describes the type of the token that was issued.
	IssuedTokenType string `json:"issued_token_type"`

	// TokenType describes the format or usage type of the issued token (e.g., Bearer).
	TokenType string `json:"token_type"`

	// ExpiresIn indicates how many seconds the token is valid from the time of issuance.
	ExpiresIn int64 `json:"expires_in"`
}

// ExchangeToken performs a token exchange with the tokenbridge service.
// It sends a POST request with the service account token and optionally an audience
// to obtain a new access token.
//
// Parameters:
//   - exchangeURL: The URL of the tokenbridge exchange endpoint.
//   - idToken: The Kubernetes service account token to exchange.
//   - audience: Optional. Specifies the intended audience for the requested token.
//
// Returns:
//   - *ExchangeResponse containing the new access token and metadata.
//   - error if the request fails or the response is invalid.
func ExchangeToken(exchangeURL, idToken string, audience string) (*ExchangeResponse, error) {
	values := url.Values{}
	values.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	values.Set("subject_token", idToken)
	values.Set("subject_token_type", "urn:ietf:params:oauth:token-type:id_token")

	if audience != "" {
		values.Set("audience", audience)
	}

	req, err := http.NewRequest("POST", exchangeURL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-store")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exchange failed: %s - %s", resp.Status, string(body))
	}

	var result ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
