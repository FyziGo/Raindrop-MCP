package api

import (
	"encoding/json"
	"fmt"

	"raindrop-mcp/types"
)

// CreateCollection creates a new collection
func (c *Client) CreateCollection(title string, parentID int, isPublic bool) (*types.Collection, error) {
	reqBody := types.CreateCollectionRequest{
		Title:  title,
		Public: isPublic,
	}
	if parentID > 0 {
		reqBody.Parent = &types.CollectionRef{ID: parentID}
	}

	respBody, err := c.makeRequest("POST", "/collection", reqBody)
	if err != nil {
		return nil, err
	}

	var resp types.SingleCollectionResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// GetCollection retrieves a single collection by ID
func (c *Client) GetCollection(id int) (*types.Collection, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/collection/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var resp types.SingleCollectionResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// UpdateCollection updates an existing collection
func (c *Client) UpdateCollection(id int, title string, isPublic *bool, parentID *int) (*types.Collection, error) {
	reqBody := types.UpdateCollectionRequest{}

	if title != "" {
		reqBody.Title = title
	}
	if isPublic != nil {
		reqBody.Public = isPublic
	}
	if parentID != nil {
		reqBody.Parent = &types.CollectionRef{ID: *parentID}
	}

	respBody, err := c.makeRequest("PUT", fmt.Sprintf("/collection/%d", id), reqBody)
	if err != nil {
		return nil, err
	}

	var resp types.SingleCollectionResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// DeleteCollection removes a collection
func (c *Client) DeleteCollection(id int) error {
	_, err := c.makeRequest("DELETE", fmt.Sprintf("/collection/%d", id), nil)
	return err
}

// MergeCollections merges collections into target
func (c *Client) MergeCollections(ids []int, targetID int) error {
	reqBody := map[string]any{
		"to":  targetID,
		"ids": ids,
	}
	_, err := c.makeRequest("PUT", "/collections/merge", reqBody)
	return err
}
