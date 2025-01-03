package models

type PageInfo struct {
	HTMLVersion       string         `json:"html_version"`
	Title             string         `json:"title"`
	HeadingCounts     map[string]int `json:"heading_counts"`
	URLs              []URLStatus    `json:"urls"`
	InternalURLsCount int            `json:"internal_urls_count"`
	ExternalURLsCount int            `json:"external_urls_count"`
	ContainsLoginForm bool           `json:"contains_login_form"`
}

type URLStatus struct {
	URL        string `json:"url"`
	HTTPStatus int    `json:"http_status"`
	Error      string `json:"error"`
}
