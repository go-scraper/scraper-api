package models

type Pagination struct {
	PageSize    int     `json:"page_size"`
	CurrentPage int     `json:"current_page"`
	TotalPages  int     `json:"total_pages"`
	PrevPage    *string `json:"prev_page,omitempty"`
	NextPage    *string `json:"next_page,omitempty"`
}

type ScrapedData struct {
	HTMLVersion       string         `json:"html_version"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	ContainsLoginForm bool           `json:"contains_login_form"`
	TotalURLs         int            `json:"total_urls"`
	InternalURLs      int            `json:"internal_urls"`
	ExternalURLs      int            `json:"external_urls"`
	Paginated         PaginatedURLs  `json:"paginated"`
}

type PaginatedURLs struct {
	InaccessibleURLs int         `json:"inaccessible_urls"`
	URLs             []URLStatus `json:"urls"`
}

type PageResponse struct {
	RequestID  string      `json:"request_id"`
	SessionId  string      `json:"session_id"`
	Pagination Pagination  `json:"pagination"`
	Scraped    ScrapedData `json:"scraped"`
}
