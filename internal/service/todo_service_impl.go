package service

import (
	"errors"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

var (
	// ErrInvalidDueDate は日付形式が不正な場合のエラーです
	ErrInvalidDueDate = errors.New("invalid due date format")
	// ErrEmptyTitle はタイトルが空の場合のエラーです
	ErrEmptyTitle = errors.New("title cannot be empty")
)

// TodoServiceImpl はTodoServiceインターフェースの実装です
type TodoServiceImpl struct {
	todoRepo repository.TodoRepository
}

// NewTodoService は新しいTodoServiceImplのインスタンスを生成します
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &TodoServiceImpl{
		todoRepo: todoRepo,
	}
}

// GetByID は指定されたIDのTODOアイテムを取得します
func (s *TodoServiceImpl) GetByID(id int) (*domain.Todo, error) {
	return s.todoRepo.GetByID(id)
}

// Create は新しいTODOアイテムを作成します
func (s *TodoServiceImpl) Create(title string, dueDate string) (*domain.Todo, error) {
	// バリデーション
	if title == "" {
		return nil, ErrEmptyTitle
	}

	// 日付文字列をtime.Time型に変換
	parsedDueDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return nil, ErrInvalidDueDate
	}

	// TODOアイテムを作成
	todo := &domain.Todo{
		Title:   title,
		DueDate: parsedDueDate,
	}

	// リポジトリに保存
	return s.todoRepo.Create(todo)
}
