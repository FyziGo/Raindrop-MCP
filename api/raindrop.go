package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"raindrop-mcp/types"
)

const baseURL = "https://api.raindrop.io/rest/v1"

// Maximum response size (10MB)
const maxResponseSize = 10 * 1024 * 1024

// Client is the Raindrop.io API client
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new Raindrop API client
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest performs an HTTP request to the Raindrop API
func (c *Client) makeRequest(method, endpoint string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CreateRaindrop creates a new bookmark
func (c *Client) CreateRaindrop(link, title string, tags []string, collectionID int) (*types.Raindrop, error) {
	reqBody := types.CreateRaindropRequest{
		Link:        link,
		Title:       title,
		Tags:        tags,
		PleaseParse: map[string]any{}, // Enable auto-parsing of metadata
	}
	if collectionID != 0 {
		reqBody.Collection = &types.CollectionRef{ID: collectionID}
	}

	respBody, err := c.makeRequest("POST", "/raindrop", reqBody)
	if err != nil {
		return nil, err
	}

	var resp types.SingleRaindropResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// GetRaindrop retrieves a single bookmark by ID
func (c *Client) GetRaindrop(id int) (*types.Raindrop, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/raindrop/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var resp types.SingleRaindropResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// UpdateRaindrop updates an existing bookmark
func (c *Client) UpdateRaindrop(id int, title, note string, tags []string, collectionID *int) (*types.Raindrop, error) {
	reqBody := types.UpdateRaindropRequest{}

	if title != "" {
		reqBody.Title = title
	}
	if note != "" {
		reqBody.Note = note
	}
	if tags != nil {
		reqBody.Tags = tags
	}
	if collectionID != nil {
		reqBody.Collection = &types.CollectionRef{ID: *collectionID}
	}

	respBody, err := c.makeRequest("PUT", fmt.Sprintf("/raindrop/%d", id), reqBody)
	if err != nil {
		return nil, err
	}

	var resp types.SingleRaindropResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Item, nil
}

// DeleteRaindrop deletes a bookmark (moves to Trash)
func (c *Client) DeleteRaindrop(id int) error {
	_, err := c.makeRequest("DELETE", fmt.Sprintf("/raindrop/%d", id), nil)
	return err
}

// SearchRaindrops searches for bookmarks
func (c *Client) SearchRaindrops(query string, collectionID, page, perPage int, tags []string) (*types.RaindropsResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Set("search", query)
	}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("perpage", strconv.Itoa(perPage))
	}

	endpoint := fmt.Sprintf("/raindrops/%d", collectionID)
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp types.RaindropsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ListCollections returns all root collections
func (c *Client) ListCollections() (*types.CollectionsResponse, error) {
	respBody, err := c.makeRequest("GET", "/collections", nil)
	if err != nil {
		return nil, err
	}

	var resp types.CollectionsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ListChildCollections returns all nested collections
func (c *Client) ListChildCollections() (*types.CollectionsResponse, error) {
	respBody, err := c.makeRequest("GET", "/collections/childrens", nil)
	if err != nil {
		return nil, err
	}

	var resp types.CollectionsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetTags returns all tags, optionally filtered by collection
func (c *Client) GetTags(collectionID int) (*types.TagsResponse, error) {
	endpoint := "/tags"
	if collectionID != 0 {
		endpoint = fmt.Sprintf("/tags/%d", collectionID)
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp types.TagsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}
