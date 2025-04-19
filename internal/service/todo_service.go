// TodoService はTODOアイテムのサービスのインターフェースです
package service

import (
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
)

// TodoService はTODOアイテムのサービスのインターフェースです
type TodoService interface {
	// GetByID は指定されたIDのTODOアイテムを取得します
	GetByID(id int) (*domain.Todo, error)

	// Create は新しいTODOアイテムを作成します
	Create(title string, dueDate string) (*domain.Todo, error)
}
