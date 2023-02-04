package cron

type Response struct {
	Items         []SearchResult `json:"items,omitempty"`
	NextPageToken string         `json:"nextPageToken,omitempty"`
	PrevPageToken string         `json:"prevPageToken,omitempty"`
}

type SearchResult struct {
	Snippet SearchResultSnippet `json:"snippet,omitempty"`
}

type SearchResultSnippet struct {
	Description string           `json:"description,omitempty"`
	PublishedAt string           `json:"publishedAt,omitempty"`
	Thumbnails  ThumbnailDetails `json:"thumbnails,omitempty"`
	Title       string           `json:"title,omitempty"`
}

type ThumbnailDetails struct {
	Default  Thumbnail `json:"default,omitempty"`
	High     Thumbnail `json:"high,omitempty"`
	Maxres   Thumbnail `json:"maxres,omitempty"`
	Medium   Thumbnail `json:"medium,omitempty"`
	Standard Thumbnail `json:"standard,omitempty"`
}

type Thumbnail struct {
	Height int64  `json:"height,omitempty"`
	Url    string `json:"url,omitempty"`
	Width  int64  `json:"width,omitempty"`
}
