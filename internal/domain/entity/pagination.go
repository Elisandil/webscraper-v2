package entity

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

type PaginatedScrapingResults struct {
	Data       []*ScrapingResult   `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

func NewPaginationRequest(page, perPage int) *PaginationRequest {

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	return &PaginationRequest{
		Page:    page,
		PerPage: perPage,
	}
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func NewPaginationResponse(page, perPage int, totalItems int64) *PaginationResponse {
	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))

	return &PaginationResponse{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}
