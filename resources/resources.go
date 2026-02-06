package resources

import (
	"context"
	"fmt"
	"strings"

	"raindrop-mcp/api"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterResources registers MCP resources for Raindrop data
func RegisterResources(server *mcp.Server, client *api.Client) {
	// Resource: All collections
	server.AddResource(&mcp.Resource{
		URI:         "raindrop://collections",
		Name:        "All Collections",
		Description: "List of all Raindrop.io collections",
		MIMEType:    "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		rootCollections, err := client.ListCollections()
		if err != nil {
			return nil, fmt.Errorf("failed to list collections: %w", err)
		}
		childCollections, err := client.ListChildCollections()
		if err != nil {
			return nil, fmt.Errorf("failed to list child collections: %w", err)
		}

		var sb strings.Builder
		sb.WriteString("# Raindrop.io Collections\n\n")

		allCollections := append(rootCollections.Items, childCollections.Items...)
		for _, c := range allCollections {
			indent := ""
			if c.Parent != nil && c.Parent.ID > 0 {
				indent = "  "
			}
			sb.WriteString(fmt.Sprintf("%s- %s (ID: %d, %d bookmarks)\n", indent, c.Title, c.FullID, c.Count))
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: sb.String(),
			}},
		}, nil
	})

	// Resource: All tags
	server.AddResource(&mcp.Resource{
		URI:         "raindrop://tags",
		Name:        "All Tags",
		Description: "List of all Raindrop.io tags",
		MIMEType:    "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		tags, err := client.GetTags(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}

		var sb strings.Builder
		sb.WriteString("# Raindrop.io Tags\n\n")

		for _, t := range tags.Items {
			sb.WriteString(fmt.Sprintf("- %s (%d bookmarks)\n", t.ID, t.Count))
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: sb.String(),
			}},
		}, nil
	})

	// Resource: User info
	server.AddResource(&mcp.Resource{
		URI:         "raindrop://user",
		Name:        "User Info",
		Description: "Current Raindrop.io user information",
		MIMEType:    "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		user, err := client.GetUser()
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		var sb strings.Builder
		sb.WriteString("# Raindrop.io User\n\n")
		sb.WriteString(fmt.Sprintf("Name: %s\n", user.FullName))
		sb.WriteString(fmt.Sprintf("Email: %s\n", user.Email))
		if user.Pro {
			sb.WriteString("Account: Pro\n")
		} else {
			sb.WriteString("Account: Free\n")
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: sb.String(),
			}},
		}, nil
	})

	// Resource Template: Bookmarks in collection
	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "raindrop://collection/{id}/bookmarks",
		Name:        "Collection Bookmarks",
		Description: "Bookmarks in a specific collection",
		MIMEType:    "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		// Parse collection ID from URI
		uri := req.Params.URI
		var collectionID int
		_, err := fmt.Sscanf(uri, "raindrop://collection/%d/bookmarks", &collectionID)
		if err != nil {
			return nil, fmt.Errorf("invalid collection URI: %w", err)
		}

		raindrops, err := client.SearchRaindrops("", collectionID, 0, 25, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get bookmarks: %w", err)
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# Bookmarks in Collection %d\n\n", collectionID))
		sb.WriteString(fmt.Sprintf("Total: %d bookmarks\n\n", raindrops.Count))

		for _, r := range raindrops.Items {
			sb.WriteString(fmt.Sprintf("- **%s**\n", r.Title))
			sb.WriteString(fmt.Sprintf("  URL: %s\n", r.Link))
			if len(r.Tags) > 0 {
				sb.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(r.Tags, ", ")))
			}
			sb.WriteString("\n")
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: sb.String(),
			}},
		}, nil
	})
}
