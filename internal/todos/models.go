package todos

import "gorm.io/gorm"

type TodoItem struct {
	gorm.Model
	Title       string
	Description string
	UserId      uint
}
