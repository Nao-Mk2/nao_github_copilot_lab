package repository

import (
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
)

// TodoRepository はTODOアイテムのリポジトリのインターフェースです
type TodoRepository interface {
	// GetByID は指定されたIDのTODOアイテムを取得します
	GetByID(id int) (*domain.Todo, error)

	// Create は新しいTODOアイテムを作成します
	Create(todo *domain.Todo) (*domain.Todo, error)
}
