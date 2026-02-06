package types

// User represents a Raindrop.io user
type User struct {
	ID         int     `json:"_id"`
	Email      string  `json:"email"`
	EmailMD5   string  `json:"email_MD5"`
	FullName   string  `json:"fullName"`
	Pro        bool    `json:"pro"`
	ProExpire  string  `json:"pro_expire,omitempty"`
	Registered string  `json:"registered"`
	Groups     []Group `json:"groups"`
	Config     Config  `json:"config"`
}

// Group represents a collection group
type Group struct {
	Title       string `json:"title"`
	Hidden      bool   `json:"hidden"`
	Sort        int    `json:"sort"`
	Collections []int  `json:"collections"`
}

// Config represents user configuration
type Config struct {
	RaindropSort string `json:"raindrops_sort"`
}

// UserResponse is the response for user endpoint
type UserResponse struct {
	Result bool `json:"result"`
	User   User `json:"user"`
}

// Highlight represents a text highlight in a raindrop
type Highlight struct {
	ID         string `json:"_id"`
	RaindropID int    `json:"raindropRef,omitempty"`
	Text       string `json:"text"`
	Note       string `json:"note,omitempty"`
	Color      string `json:"color,omitempty"`
	Created    string `json:"created,omitempty"`
	LastUpdate string `json:"lastUpdate,omitempty"`
}

// HighlightsResponse is the response for highlights list
type HighlightsResponse struct {
	Result bool        `json:"result"`
	Items  []Highlight `json:"items"`
}

// SingleHighlightResponse is the response for single highlight operations
type SingleHighlightResponse struct {
	Result bool      `json:"result"`
	Item   Highlight `json:"item"`
}

// Filter represents a filter in a collection
type Filter struct {
	Type       string       `json:"type,omitempty"`
	Broken     string       `json:"broken,omitempty"`
	Duplicates string       `json:"duplicates,omitempty"`
	Tags       []string     `json:"tags,omitempty"`
	Types      []TypeFilter `json:"types,omitempty"`
}

// TypeFilter represents a type filter
type TypeFilter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// FiltersResponse is the response for filters
type FiltersResponse struct {
	Result bool     `json:"result"`
	Items  []Filter `json:"items,omitempty"`
	Broken struct {
		Count int `json:"count"`
	} `json:"broken,omitempty"`
	Duplicates struct {
		Count int `json:"count"`
	} `json:"duplicates,omitempty"`
	Tags  []TagWithCount `json:"tags,omitempty"`
	Types []TypeFilter   `json:"types,omitempty"`
}

// TagWithCount represents a tag with count for filters
type TagWithCount struct {
	ID    string `json:"_id"`
	Count int    `json:"count"`
}

// CreateCollectionRequest is the request for creating a collection
type CreateCollectionRequest struct {
	Title  string         `json:"title"`
	Sort   int            `json:"sort,omitempty"`
	Public bool           `json:"public,omitempty"`
	Parent *CollectionRef `json:"parent,omitempty"`
	View   string         `json:"view,omitempty"`
	Cover  []string       `json:"cover,omitempty"`
}

// UpdateCollectionRequest is the request for updating a collection
type UpdateCollectionRequest struct {
	Title    string         `json:"title,omitempty"`
	Sort     int            `json:"sort,omitempty"`
	Public   *bool          `json:"public,omitempty"`
	Parent   *CollectionRef `json:"parent,omitempty"`
	View     string         `json:"view,omitempty"`
	Cover    []string       `json:"cover,omitempty"`
	Expanded bool           `json:"expanded,omitempty"`
}

// SingleCollectionResponse is the response for single collection operations
type SingleCollectionResponse struct {
	Result bool       `json:"result"`
	Item   Collection `json:"item"`
}

// BulkRaindropsRequest is the request for bulk operations
type BulkRaindropsRequest struct {
	Items []CreateRaindropRequest `json:"items"`
}

// BulkUpdateRequest is the request for bulk update
type BulkUpdateRequest struct {
	IDs        []int    `json:"ids"`
	Tags       []string `json:"tags,omitempty"`
	Collection int      `json:"collection,omitempty"`
	Important  *bool    `json:"important,omitempty"`
}

// BulkDeleteRequest is the request for bulk delete
type BulkDeleteRequest struct {
	IDs []int `json:"ids"`
}

// BulkResponse is the response for bulk operations
type BulkResponse struct {
	Result   bool `json:"result"`
	Modified int  `json:"modified,omitempty"`
}

// CreateHighlightRequest is the request for creating a highlight
type CreateHighlightRequest struct {
	RaindropID int    `json:"raindrop"`
	Text       string `json:"text"`
	Note       string `json:"note,omitempty"`
	Color      string `json:"color,omitempty"`
}

// RenameTagRequest is the request for renaming a tag
type RenameTagRequest struct {
	OldName string `json:"old"`
	NewName string `json:"new"`
}

// MergeTagsRequest is the request for merging tags
type MergeTagsRequest struct {
	Tags []string `json:"tags"`
}

// SuggestResponse is the response for suggestions
type SuggestResponse struct {
	Result bool     `json:"result"`
	Items  []string `json:"items,omitempty"`
}
