package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"
)

const (
	authURL  = "https://raindrop.io/oauth/authorize"
	tokenURL = "https://raindrop.io/oauth/access_token"
)

// OAuthConfig contains OAuth2 configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// tokenResponse represents the OAuth token response
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error,omitempty"`
}

// StartOAuthFlow initiates the OAuth2 authorization flow
// Opens browser for user authorization and waits for callback
func StartOAuthFlow(ctx context.Context, config *OAuthConfig) (*TokenData, error) {
	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)
	config.RedirectURI = redirectURI

	// Channel to receive authorization code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Setup callback handler
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		errParam := r.URL.Query().Get("error")

		if errParam != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body><h1>Authorization Failed</h1><p>Error: %s</p></body></html>", errParam)
			errChan <- fmt.Errorf("authorization denied: %s", errParam)
			return
		}

		if code == "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, "<html><body><h1>Error</h1><p>No authorization code received</p></body></html>")
			errChan <- fmt.Errorf("no authorization code received")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body>
			<h1>✅ Authorization Successful!</h1>
			<p>You can close this window and return to your application.</p>
			<script>setTimeout(function(){window.close();}, 2000);</script>
		</body></html>`)
		codeChan <- code
	})

	server := &http.Server{Handler: mux}

	// Start server in background
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("callback server error: %w", err)
		}
	}()

	// Build authorization URL
	authParams := url.Values{}
	authParams.Set("client_id", config.ClientID)
	authParams.Set("redirect_uri", redirectURI)
	authParams.Set("response_type", "code")
	authURLFull := authURL + "?" + authParams.Encode()

	// Open browser
	fmt.Printf("Opening browser for authorization...\n")
	fmt.Printf("If browser doesn't open, visit: %s\n", authURLFull)
	if err := openBrowser(authURLFull); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}

	// Wait for code or error
	var code string
	select {
	case code = <-codeChan:
		// Success
	case err := <-errChan:
		server.Shutdown(ctx)
		return nil, err
	case <-time.After(5 * time.Minute):
		server.Shutdown(ctx)
		return nil, fmt.Errorf("authorization timeout (5 minutes)")
	case <-ctx.Done():
		server.Shutdown(ctx)
		return nil, ctx.Err()
	}

	// Shutdown callback server
	server.Shutdown(ctx)

	// Exchange code for token
	return exchangeCodeForToken(config, code)
}

// exchangeCodeForToken exchanges authorization code for access token
func exchangeCodeForToken(config *OAuthConfig, code string) (*TokenData, error) {
	reqBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret,
		"redirect_uri":  config.RedirectURI,
	}

	return makeTokenRequest(reqBody)
}

// RefreshToken refreshes an expired access token
func RefreshToken(config *OAuthConfig, refreshToken string) (*TokenData, error) {
	reqBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret,
	}

	return makeTokenRequest(reqBody)
}

// makeTokenRequest makes a token request to Raindrop API
func makeTokenRequest(params map[string]string) (*TokenData, error) {
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(tokenURL, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB limit
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token error: %s", tokenResp.Error)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response")
	}

	return &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + tokenResp.ExpiresIn,
		TokenType:    tokenResp.TokenType,
	}, nil
}

// openBrowser opens URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // Linux and others
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
