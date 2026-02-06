package api

import (
	"encoding/json"
	"fmt"

	"raindrop-mcp/types"
)

// RenameTag renames a tag in a collection (0 for all)
func (c *Client) RenameTag(collectionID int, oldName, newName string) error {
	reqBody := []types.RenameTagRequest{
		{OldName: oldName, NewName: newName},
	}

	endpoint := "/tags"
	if collectionID > 0 {
		endpoint = fmt.Sprintf("/tags/%d", collectionID)
	}

	_, err := c.makeRequest("PUT", endpoint, reqBody)
	return err
}

// DeleteTags deletes tags in a collection (0 for all)
func (c *Client) DeleteTags(collectionID int, tags []string) error {
	reqBody := map[string][]string{
		"tags": tags,
	}

	endpoint := "/tags"
	if collectionID > 0 {
		endpoint = fmt.Sprintf("/tags/%d", collectionID)
	}

	_, err := c.makeRequest("DELETE", endpoint, reqBody)
	return err
}

// MergeTags merges multiple tags into one
func (c *Client) MergeTags(collectionID int, tags []string) error {
	reqBody := types.MergeTagsRequest{
		Tags: tags,
	}

	endpoint := "/tags/merge"
	if collectionID > 0 {
		endpoint = fmt.Sprintf("/tags/%d/merge", collectionID)
	}

	_, err := c.makeRequest("PUT", endpoint, reqBody)
	return err
}

// GetHighlights gets all highlights for a raindrop
func (c *Client) GetHighlights(raindropID int) (*types.HighlightsResponse, error) {
	endpoint := fmt.Sprintf("/raindrop/%d/highlights", raindropID)

	// If raindropID is 0, get all highlights
	if raindropID == 0 {
		endpoint = "/highlights"
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp types.HighlightsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// CreateHighlight creates a new highlight
func (c *Client) CreateHighlight(raindropID int, text, note, color string) (*types.Highlight, error) {
	reqBody := types.CreateHighlightRequest{
		RaindropID: raindropID,
		Text:       text,
		Note:       note,
		Color:      color,
	}

	respBody, err := c.makeRequest("POST", "/highlight", reqBody)
	if err != nil {
		return nil, err
	}

	var resp types.SingleHighlightResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// DeleteHighlight deletes a highlight
func (c *Client) DeleteHighlight(raindropID int, highlightID string) error {
	_, err := c.makeRequest("DELETE", fmt.Sprintf("/raindrop/%d/highlight/%s", raindropID, highlightID), nil)
	return err
}

// GetFilters gets filters for a collection
func (c *Client) GetFilters(collectionID int) (*types.FiltersResponse, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/filters/%d", collectionID), nil)
	if err != nil {
		return nil, err
	}

	var resp types.FiltersResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetUser gets current user info
func (c *Client) GetUser() (*types.User, error) {
	respBody, err := c.makeRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var resp types.UserResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.User, nil
}

// SuggestTags suggests tags for a URL
func (c *Client) SuggestTags(url string) ([]string, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/tags/suggest?url=%s", url), nil)
	if err != nil {
		return nil, err
	}

	var resp types.SuggestResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Items, nil
}
