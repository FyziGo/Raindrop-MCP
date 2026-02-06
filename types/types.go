package types

// Raindrop represents a bookmark in Raindrop.io
type Raindrop struct {
	ID         int        `json:"_id"`
	Collection Collection `json:"collection"`
	Cover      string     `json:"cover"`
	Created    string     `json:"created"`
	Domain     string     `json:"domain"`
	Excerpt    string     `json:"excerpt"`
	LastUpdate string     `json:"lastUpdate"`
	Link       string     `json:"link"`
	Media      []Media    `json:"media"`
	Tags       []string   `json:"tags"`
	Title      string     `json:"title"`
	Type       string     `json:"type"`
	Note       string     `json:"note"`
	Important  bool       `json:"important"`
}

// Collection represents a Raindrop.io collection
type Collection struct {
	ID         int    `json:"$id,omitempty"`
	FullID     int    `json:"_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Count      int    `json:"count,omitempty"`
	Cover      []string `json:"cover,omitempty"`
	Color      string `json:"color,omitempty"`
	Created    string `json:"created,omitempty"`
	LastUpdate string `json:"lastUpdate,omitempty"`
	Public     bool   `json:"public,omitempty"`
	View       string `json:"view,omitempty"`
	Parent     *Parent `json:"parent,omitempty"`
}

// Parent represents parent collection reference
type Parent struct {
	ID int `json:"$id"`
}

// Media represents media item in a raindrop
type Media struct {
	Link string `json:"link"`
	Type string `json:"type,omitempty"`
}

// Tag represents a tag with count
type Tag struct {
	ID    string `json:"_id"`
	Count int    `json:"count"`
}

// API Response types

// SingleRaindropResponse is the response for single raindrop operations
type SingleRaindropResponse struct {
	Result bool     `json:"result"`
	Item   Raindrop `json:"item"`
}

// RaindropsResponse is the response for multiple raindrops
type RaindropsResponse struct {
	Result bool       `json:"result"`
	Items  []Raindrop `json:"items"`
	Count  int        `json:"count"`
}

// CollectionsResponse is the response for collections list
type CollectionsResponse struct {
	Result bool         `json:"result"`
	Items  []Collection `json:"items"`
}

// TagsResponse is the response for tags list
type TagsResponse struct {
	Result bool  `json:"result"`
	Items  []Tag `json:"items"`
}

// CreateRaindropRequest is the request body for creating a raindrop
type CreateRaindropRequest struct {
	Link        string            `json:"link"`
	Title       string            `json:"title,omitempty"`
	Excerpt     string            `json:"excerpt,omitempty"`
	Note        string            `json:"note,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Important   bool              `json:"important,omitempty"`
	Collection  *CollectionRef    `json:"collection,omitempty"`
	PleaseParse map[string]any    `json:"pleaseParse,omitempty"`
}

// UpdateRaindropRequest is the request body for updating a raindrop
type UpdateRaindropRequest struct {
	Link       string         `json:"link,omitempty"`
	Title      string         `json:"title,omitempty"`
	Excerpt    string         `json:"excerpt,omitempty"`
	Note       string         `json:"note,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	Important  *bool          `json:"important,omitempty"`
	Collection *CollectionRef `json:"collection,omitempty"`
}

// CollectionRef is used to reference a collection by ID
type CollectionRef struct {
	ID int `json:"$id"`
}
