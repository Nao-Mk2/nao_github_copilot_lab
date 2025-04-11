package model

import (
	"errors"
	"time"
)

var (
	// ErrEmptyTitle はタイトルが空の場合のエラー
	ErrEmptyTitle = errors.New("タイトルは空にできません")
	// ErrInvalidDueDate は期限が現在時刻より前の場合のエラー
	ErrInvalidDueDate = errors.New("期限は現在時刻より後に設定してください")
)

// Todo はTODOアイテムを表す構造体
type Todo struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	DueDate time.Time `json:"due_date"`
}

// NewTodo は新しいTodoインスタンスを作成します
func NewTodo(title string, dueDate time.Time) (*Todo, error) {
	todo := &Todo{
		Title:   title,
		DueDate: dueDate,
	}

	if err := todo.Validate(); err != nil {
		return nil, err
	}

	return todo, nil
}

// Validate はTodoオブジェクトのバリデーションを行います
func (t *Todo) Validate() error {
	if t.Title == "" {
		return ErrEmptyTitle
	}

	if t.DueDate.Before(time.Now()) {
		return ErrInvalidDueDate
	}

	return nil
}
