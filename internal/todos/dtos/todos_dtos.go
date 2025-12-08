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
