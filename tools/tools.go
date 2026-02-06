package tools

import (
	"context"
	"fmt"
	"strings"

	"raindrop-mcp/api"
	"raindrop-mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Output type for tools that return text
type TextOutput struct {
	Text string `json:"text"`
}

// RegisterTools registers all Raindrop tools with the MCP server
func RegisterTools(server *mcp.Server, client *api.Client) {
	// create-bookmark
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-bookmark",
		Description: "Create a new bookmark in Raindrop.io",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateBookmarkInput) (*mcp.CallToolResult, TextOutput, error) {
		raindrop, err := client.CreateRaindrop(input.URL, input.Title, input.Tags, input.Collection)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to create bookmark: %w", err)
		}
		return nil, TextOutput{Text: formatRaindrop(raindrop)}, nil
	})

	// get-bookmark
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-bookmark",
		Description: "Get a bookmark by its ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetBookmarkInput) (*mcp.CallToolResult, TextOutput, error) {
		raindrop, err := client.GetRaindrop(input.ID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get bookmark: %w", err)
		}
		return nil, TextOutput{Text: formatRaindrop(raindrop)}, nil
	})

	// update-bookmark
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update-bookmark",
		Description: "Update an existing bookmark",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateBookmarkInput) (*mcp.CallToolResult, TextOutput, error) {
		var collectionPtr *int
		if input.Collection != 0 {
			collectionPtr = &input.Collection
		}
		raindrop, err := client.UpdateRaindrop(input.ID, input.Title, input.Note, input.Tags, collectionPtr)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to update bookmark: %w", err)
		}
		return nil, TextOutput{Text: formatRaindrop(raindrop)}, nil
	})

	// delete-bookmark
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-bookmark",
		Description: "Delete a bookmark (moves to Trash)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteBookmarkInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.DeleteRaindrop(input.ID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to delete bookmark: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Bookmark %d deleted successfully", input.ID)}, nil
	})

	// search-bookmarks
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search-bookmarks",
		Description: "Search through your Raindrop.io bookmarks",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchBookmarksInput) (*mcp.CallToolResult, TextOutput, error) {
		result, err := client.SearchRaindrops(input.Query, input.Collection, input.Page, input.PerPage, input.Tags)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to search bookmarks: %w", err)
		}
		return nil, TextOutput{Text: formatRaindrops(result)}, nil
	})

	// list-collections
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-collections",
		Description: "List all your Raindrop.io collections",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, TextOutput, error) {
		rootCollections, err := client.ListCollections()
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to list collections: %w", err)
		}
		childCollections, err := client.ListChildCollections()
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to list child collections: %w", err)
		}

		allCollections := append(rootCollections.Items, childCollections.Items...)
		return nil, TextOutput{Text: formatCollections(allCollections)}, nil
	})

	// list-tags
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-tags",
		Description: "List all tags in your Raindrop.io account",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListTagsInput) (*mcp.CallToolResult, TextOutput, error) {
		tagsResp, err := client.GetTags(input.Collection)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to list tags: %w", err)
		}
		return nil, TextOutput{Text: formatTags(tagsResp.Items)}, nil
	})
}

// Input types

type CreateBookmarkInput struct {
	URL        string   `json:"url" jsonschema:"description=URL to bookmark,required"`
	Title      string   `json:"title,omitempty" jsonschema:"description=Title for the bookmark"`
	Tags       []string `json:"tags,omitempty" jsonschema:"description=Tags for the bookmark"`
	Collection int      `json:"collection,omitempty" jsonschema:"description=Collection ID to save to (0 for Unsorted)"`
}

type GetBookmarkInput struct {
	ID int `json:"id" jsonschema:"description=Bookmark ID,required"`
}

type UpdateBookmarkInput struct {
	ID         int      `json:"id" jsonschema:"description=Bookmark ID to update,required"`
	Title      string   `json:"title,omitempty" jsonschema:"description=New title"`
	Note       string   `json:"note,omitempty" jsonschema:"description=Note/description"`
	Tags       []string `json:"tags,omitempty" jsonschema:"description=New tags (replaces existing)"`
	Collection int      `json:"collection,omitempty" jsonschema:"description=Move to collection ID"`
}

type DeleteBookmarkInput struct {
	ID int `json:"id" jsonschema:"description=Bookmark ID to delete,required"`
}

type SearchBookmarksInput struct {
	Query      string   `json:"query" jsonschema:"description=Search query,required"`
	Collection int      `json:"collection,omitempty" jsonschema:"description=Collection ID (0 for all)"`
	Tags       []string `json:"tags,omitempty" jsonschema:"description=Filter by tags"`
	Page       int      `json:"page,omitempty" jsonschema:"description=Page number (0-based)"`
	PerPage    int      `json:"perpage,omitempty" jsonschema:"description=Items per page (max 50)"`
}

type ListTagsInput struct {
	Collection int `json:"collection,omitempty" jsonschema:"description=Collection ID (0 for all tags)"`
}

// Formatting helpers

func formatRaindrop(r *types.Raindrop) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s**\n", r.Title))
	sb.WriteString(fmt.Sprintf("ID: %d\n", r.ID))
	sb.WriteString(fmt.Sprintf("URL: %s\n", r.Link))
	if len(r.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(r.Tags, ", ")))
	}
	if r.Excerpt != "" {
		sb.WriteString(fmt.Sprintf("Excerpt: %s\n", r.Excerpt))
	}
	if r.Note != "" {
		sb.WriteString(fmt.Sprintf("Note: %s\n", r.Note))
	}
	sb.WriteString(fmt.Sprintf("Created: %s\n", r.Created))
	return sb.String()
}

func formatRaindrops(resp *types.RaindropsResponse) string {
	if len(resp.Items) == 0 {
		return "No bookmarks found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d bookmarks:\n\n", resp.Count))

	for i, r := range resp.Items {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, r.Title))
		sb.WriteString(fmt.Sprintf("   ID: %d | URL: %s\n", r.ID, r.Link))
		if len(r.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("   Tags: %s\n", strings.Join(r.Tags, ", ")))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatCollections(collections []types.Collection) string {
	if len(collections) == 0 {
		return "No collections found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d collections:\n\n", len(collections)))

	for _, c := range collections {
		indent := ""
		if c.Parent != nil && c.Parent.ID > 0 {
			indent = "  └─ "
		}
		sb.WriteString(fmt.Sprintf("%s**%s** (ID: %d)\n", indent, c.Title, c.FullID))
		sb.WriteString(fmt.Sprintf("%s   %d bookmarks\n\n", indent, c.Count))
	}

	return sb.String()
}

func formatTags(tags []types.Tag) string {
	if len(tags) == 0 {
		return "No tags found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d tags:\n\n", len(tags)))

	for _, t := range tags {
		sb.WriteString(fmt.Sprintf("- **%s** (%d)\n", t.ID, t.Count))
	}

	return sb.String()
}
