package domain

import (
	"time"
)

// Todo はTODOアイテムを表す構造体です
type Todo struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	DueDate time.Time `json:"due_date"`
}
