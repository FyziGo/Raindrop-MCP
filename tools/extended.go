package tools

import (
	"context"
	"fmt"
	"strings"

	"raindrop-mcp/api"
	"raindrop-mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterExtendedTools registers additional Raindrop tools
func RegisterExtendedTools(server *mcp.Server, client *api.Client) {
	// --- Collections ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-collection",
		Description: "Create a new collection in Raindrop.io",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateCollectionInput) (*mcp.CallToolResult, TextOutput, error) {
		collection, err := client.CreateCollection(input.Title, input.Parent, input.Public)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to create collection: %w", err)
		}
		return nil, TextOutput{Text: formatCollection(collection)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-collection",
		Description: "Get a collection by its ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetCollectionInput) (*mcp.CallToolResult, TextOutput, error) {
		collection, err := client.GetCollection(input.ID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get collection: %w", err)
		}
		return nil, TextOutput{Text: formatCollection(collection)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update-collection",
		Description: "Update an existing collection",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input UpdateCollectionInput) (*mcp.CallToolResult, TextOutput, error) {
		var publicPtr *bool
		if input.Public != nil {
			publicPtr = input.Public
		}
		var parentPtr *int
		if input.Parent != 0 {
			parentPtr = &input.Parent
		}
		collection, err := client.UpdateCollection(input.ID, input.Title, publicPtr, parentPtr)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to update collection: %w", err)
		}
		return nil, TextOutput{Text: formatCollection(collection)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-collection",
		Description: "Delete a collection",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteCollectionInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.DeleteCollection(input.ID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to delete collection: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Collection %d deleted successfully", input.ID)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "merge-collections",
		Description: "Merge multiple collections into one",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MergeCollectionsInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.MergeCollections(input.IDs, input.TargetID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to merge collections: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Merged %d collections into collection %d", len(input.IDs), input.TargetID)}, nil
	})

	// --- Tags ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rename-tag",
		Description: "Rename a tag",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input RenameTagInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.RenameTag(input.Collection, input.OldName, input.NewName)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to rename tag: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Tag '%s' renamed to '%s'", input.OldName, input.NewName)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-tags",
		Description: "Delete one or more tags",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteTagsInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.DeleteTags(input.Collection, input.Tags)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to delete tags: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Deleted %d tags", len(input.Tags))}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "merge-tags",
		Description: "Merge multiple tags into one (first tag becomes the merged name)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MergeTagsInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.MergeTags(input.Collection, input.Tags)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to merge tags: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Merged tags into '%s'", input.Tags[0])}, nil
	})

	// --- Highlights ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-highlights",
		Description: "Get highlights from a bookmark or all highlights",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetHighlightsInput) (*mcp.CallToolResult, TextOutput, error) {
		highlights, err := client.GetHighlights(input.RaindropID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get highlights: %w", err)
		}
		return nil, TextOutput{Text: formatHighlights(highlights.Items)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-highlight",
		Description: "Create a new highlight in a bookmark",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreateHighlightInput) (*mcp.CallToolResult, TextOutput, error) {
		highlight, err := client.CreateHighlight(input.RaindropID, input.Text, input.Note, input.Color)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to create highlight: %w", err)
		}
		return nil, TextOutput{Text: formatHighlight(highlight)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-highlight",
		Description: "Delete a highlight",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteHighlightInput) (*mcp.CallToolResult, TextOutput, error) {
		err := client.DeleteHighlight(input.RaindropID, input.HighlightID)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to delete highlight: %w", err)
		}
		return nil, TextOutput{Text: fmt.Sprintf("Highlight deleted from bookmark %d", input.RaindropID)}, nil
	})

	// --- Filters ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-filters",
		Description: "Get available filters for a collection",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetFiltersInput) (*mcp.CallToolResult, TextOutput, error) {
		filters, err := client.GetFilters(input.Collection)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get filters: %w", err)
		}
		return nil, TextOutput{Text: formatFilters(filters)}, nil
	})

	// --- User ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-user",
		Description: "Get current user information",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, TextOutput, error) {
		user, err := client.GetUser()
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get user: %w", err)
		}
		return nil, TextOutput{Text: formatUser(user)}, nil
	})

	// --- Suggestions ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "suggest-tags",
		Description: "Get tag suggestions for a URL",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SuggestTagsInput) (*mcp.CallToolResult, TextOutput, error) {
		tags, err := client.SuggestTags(input.URL)
		if err != nil {
			return nil, TextOutput{}, fmt.Errorf("failed to get suggestions: %w", err)
		}
		if len(tags) == 0 {
			return nil, TextOutput{Text: "No tag suggestions available for this URL"}, nil
		}
		return nil, TextOutput{Text: fmt.Sprintf("Suggested tags: %s", strings.Join(tags, ", "))}, nil
	})
}

// Input types for extended tools

type CreateCollectionInput struct {
	Title  string `json:"title" jsonschema:"description=Collection title,required"`
	Parent int    `json:"parent,omitempty" jsonschema:"description=Parent collection ID (0 for root)"`
	Public bool   `json:"public,omitempty" jsonschema:"description=Make collection public"`
}

type GetCollectionInput struct {
	ID int `json:"id" jsonschema:"description=Collection ID,required"`
}

type UpdateCollectionInput struct {
	ID     int    `json:"id" jsonschema:"description=Collection ID,required"`
	Title  string `json:"title,omitempty" jsonschema:"description=New title"`
	Parent int    `json:"parent,omitempty" jsonschema:"description=New parent collection ID"`
	Public *bool  `json:"public,omitempty" jsonschema:"description=Make public or private"`
}

type DeleteCollectionInput struct {
	ID int `json:"id" jsonschema:"description=Collection ID to delete,required"`
}

type MergeCollectionsInput struct {
	IDs      []int `json:"ids" jsonschema:"description=Collection IDs to merge,required"`
	TargetID int   `json:"target_id" jsonschema:"description=Target collection ID,required"`
}

type RenameTagInput struct {
	Collection int    `json:"collection,omitempty" jsonschema:"description=Collection ID (0 for all)"`
	OldName    string `json:"old_name" jsonschema:"description=Current tag name,required"`
	NewName    string `json:"new_name" jsonschema:"description=New tag name,required"`
}

type DeleteTagsInput struct {
	Collection int      `json:"collection,omitempty" jsonschema:"description=Collection ID (0 for all)"`
	Tags       []string `json:"tags" jsonschema:"description=Tag names to delete,required"`
}

type MergeTagsInput struct {
	Collection int      `json:"collection,omitempty" jsonschema:"description=Collection ID (0 for all)"`
	Tags       []string `json:"tags" jsonschema:"description=Tags to merge (first becomes target),required"`
}

type GetHighlightsInput struct {
	RaindropID int `json:"raindrop_id,omitempty" jsonschema:"description=Bookmark ID (0 for all highlights)"`
}

type CreateHighlightInput struct {
	RaindropID int    `json:"raindrop_id" jsonschema:"description=Bookmark ID,required"`
	Text       string `json:"text" jsonschema:"description=Highlighted text,required"`
	Note       string `json:"note,omitempty" jsonschema:"description=Note about the highlight"`
	Color      string `json:"color,omitempty" jsonschema:"description=Highlight color (blue, red, yellow, green)"`
}

type DeleteHighlightInput struct {
	RaindropID  int    `json:"raindrop_id" jsonschema:"description=Bookmark ID,required"`
	HighlightID string `json:"highlight_id" jsonschema:"description=Highlight ID,required"`
}

type GetFiltersInput struct {
	Collection int `json:"collection" jsonschema:"description=Collection ID,required"`
}

type SuggestTagsInput struct {
	URL string `json:"url" jsonschema:"description=URL to get tag suggestions for,required"`
}

// Formatting helpers

func formatCollection(c *types.Collection) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** (ID: %d)\n", c.Title, c.FullID))
	sb.WriteString(fmt.Sprintf("Bookmarks: %d\n", c.Count))
	if c.Public {
		sb.WriteString("Public: Yes\n")
	}
	if c.Parent != nil && c.Parent.ID > 0 {
		sb.WriteString(fmt.Sprintf("Parent: %d\n", c.Parent.ID))
	}
	return sb.String()
}

func formatHighlight(h *types.Highlight) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**Highlight** (ID: %s)\n", h.ID))
	sb.WriteString(fmt.Sprintf("Text: %s\n", h.Text))
	if h.Note != "" {
		sb.WriteString(fmt.Sprintf("Note: %s\n", h.Note))
	}
	if h.Color != "" {
		sb.WriteString(fmt.Sprintf("Color: %s\n", h.Color))
	}
	return sb.String()
}

func formatHighlights(highlights []types.Highlight) string {
	if len(highlights) == 0 {
		return "No highlights found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d highlights:\n\n", len(highlights)))

	for i, h := range highlights {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, truncate(h.Text, 100)))
		if h.Note != "" {
			sb.WriteString(fmt.Sprintf("   Note: %s\n", h.Note))
		}
		sb.WriteString(fmt.Sprintf("   ID: %s\n\n", h.ID))
	}

	return sb.String()
}

func formatFilters(f *types.FiltersResponse) string {
	var sb strings.Builder
	sb.WriteString("**Filters:**\n\n")

	if f.Broken.Count > 0 {
		sb.WriteString(fmt.Sprintf("Broken links: %d\n", f.Broken.Count))
	}
	if f.Duplicates.Count > 0 {
		sb.WriteString(fmt.Sprintf("Duplicates: %d\n", f.Duplicates.Count))
	}

	if len(f.Tags) > 0 {
		sb.WriteString("\n**Tags:**\n")
		for _, t := range f.Tags {
			sb.WriteString(fmt.Sprintf("- %s (%d)\n", t.ID, t.Count))
		}
	}

	if len(f.Types) > 0 {
		sb.WriteString("\n**Types:**\n")
		for _, t := range f.Types {
			sb.WriteString(fmt.Sprintf("- %s (%d)\n", t.Name, t.Count))
		}
	}

	return sb.String()
}

func formatUser(u *types.User) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s**\n", u.FullName))
	sb.WriteString(fmt.Sprintf("Email: %s\n", u.Email))
	sb.WriteString(fmt.Sprintf("ID: %d\n", u.ID))
	if u.Pro {
		sb.WriteString("Account: Pro\n")
		if u.ProExpire != "" {
			sb.WriteString(fmt.Sprintf("Pro expires: %s\n", u.ProExpire))
		}
	} else {
		sb.WriteString("Account: Free\n")
	}
	sb.WriteString(fmt.Sprintf("Registered: %s\n", u.Registered))
	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
