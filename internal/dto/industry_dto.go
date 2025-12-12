package dto

// IndustryResponse represents the response structure for industry data
type IndustryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// IndustriesResponse represents the response structure for multiple industries
type IndustriesResponse struct {
	Industries []*IndustryResponse `json:"industries"`
	Total      int                 `json:"total"`
}
