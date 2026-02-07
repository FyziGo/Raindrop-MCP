package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TokenData represents OAuth2 token data
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

// IsExpired checks if the token is expired or will expire soon
func (t *TokenData) IsExpired() bool {
	// Refresh 5 minutes before expiry
	return time.Now().Unix() >= t.ExpiresAt-300
}

// IsValid checks if token data is valid
func (t *TokenData) IsValid() bool {
	return t.AccessToken != "" && t.RefreshToken != ""
}

// getTokenPath returns the path to token storage file
func getTokenPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".raindrop-mcp")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "token.json"), nil
}

// SaveToken saves token data to file
func SaveToken(token *TokenData) error {
	tokenPath, err := getTokenPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads token data from file
func LoadToken() (*TokenData, error) {
	tokenPath, err := getTokenPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No token saved
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

// DeleteToken removes the saved token
func DeleteToken() error {
	tokenPath, err := getTokenPath()
	if err != nil {
		return err
	}

	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}
