package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"raindrop-mcp/api"
	"raindrop-mcp/auth"
	"raindrop-mcp/resources"
	"raindrop-mcp/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Get API token from available sources
	token, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Create Raindrop API client
	client := api.NewClient(token)

	// Create MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "raindrop-mcp",
			Version: "2.0.0",
		},
		nil,
	)

	// Register all tools
	tools.RegisterTools(server, client)
	tools.RegisterExtendedTools(server, client)

	// Register resources
	resources.RegisterResources(server, client)

	// Run server on stdio transport
	fmt.Fprintln(os.Stderr, "Raindrop MCP Server v2.0.0 starting...")
	fmt.Fprintln(os.Stderr, "Loaded 20 tools, 4 resources")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// getAccessToken retrieves access token from available sources
// Priority:
// 1. RAINDROP_TOKEN environment variable (test token)
// 2. OAuth2 flow using CLIENT_ID and CLIENT_SECRET
// 3. Previously saved OAuth token
func getAccessToken() (string, error) {
	// Priority 1: Direct token from environment
	if token := os.Getenv("RAINDROP_TOKEN"); token != "" {
		fmt.Fprintln(os.Stderr, "Using token from RAINDROP_TOKEN environment variable")
		return token, nil
	}

	// Get OAuth credentials
	clientID := os.Getenv("RAINDROP_CLIENT_ID")
	clientSecret := os.Getenv("RAINDROP_CLIENT_SECRET")

	// Priority 2/3: OAuth flow
	if clientID != "" && clientSecret != "" {
		return getOAuthToken(clientID, clientSecret)
	}

	return "", fmt.Errorf("no authentication configured. Set RAINDROP_TOKEN or RAINDROP_CLIENT_ID + RAINDROP_CLIENT_SECRET")
}

// getOAuthToken handles OAuth token retrieval and refresh
func getOAuthToken(clientID, clientSecret string) (string, error) {
	config := &auth.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Try to load existing token
	savedToken, err := auth.LoadToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load saved token: %v\n", err)
	}

	// If we have a valid saved token
	if savedToken != nil && savedToken.IsValid() {
		// Check if token needs refresh
		if savedToken.IsExpired() {
			fmt.Fprintln(os.Stderr, "Access token expired, refreshing...")
			newToken, err := auth.RefreshToken(config, savedToken.RefreshToken)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to refresh token: %v. Starting new OAuth flow...\n", err)
				// Fall through to new OAuth flow
			} else {
				if err := auth.SaveToken(newToken); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to save refreshed token: %v\n", err)
				}
				fmt.Fprintln(os.Stderr, "Token refreshed successfully")
				return newToken.AccessToken, nil
			}
		} else {
			fmt.Fprintln(os.Stderr, "Using saved OAuth token")
			return savedToken.AccessToken, nil
		}
	}

	// Start new OAuth flow
	fmt.Fprintln(os.Stderr, "Starting OAuth authorization flow...")
	newToken, err := auth.StartOAuthFlow(context.Background(), config)
	if err != nil {
		return "", fmt.Errorf("OAuth flow failed: %w", err)
	}

	// Save the new token
	if err := auth.SaveToken(newToken); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save token: %v\n", err)
	}

	fmt.Fprintln(os.Stderr, "OAuth authorization successful")
	return newToken.AccessToken, nil
}
