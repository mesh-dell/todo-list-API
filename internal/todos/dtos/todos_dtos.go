package dtos

type TodoItemRequestDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoItemResponseDto struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoItemsPaginatedResponseDto struct {
	Data  []TodoItemResponseDto `json:"data"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
	Total int                   `json:"total"`
}

type QueryParams struct {
	Page        int
	Limit       int
	SearchQuery string
	SortBy      string
	Order       string
}
