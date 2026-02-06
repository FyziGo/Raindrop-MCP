package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"raindrop-mcp/api"
	"raindrop-mcp/resources"
	"raindrop-mcp/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Get API token from environment
	token := os.Getenv("RAINDROP_TOKEN")
	if token == "" {
		log.Fatal("RAINDROP_TOKEN environment variable is not set")
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
