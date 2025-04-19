package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

// MockTodoRepository はテスト用のモックリポジトリです
type MockTodoRepository struct {
	todos  map[int]*domain.Todo
	nextID int
}

// GetByIDメソッドの実装
func (m *MockTodoRepository) GetByID(id int) (*domain.Todo, error) {
	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// Createメソッドの実装
func (m *MockTodoRepository) Create(todo *domain.Todo) (*domain.Todo, error) {
	todo.ID = m.nextID
	m.nextID++
	m.todos[todo.ID] = todo
	return todo, nil
}

// NewMockTodoRepository はモックリポジトリのインスタンスを生成します
func NewMockTodoRepository() repository.TodoRepository {
	return &MockTodoRepository{
		todos:  make(map[int]*domain.Todo),
		nextID: 1,
	}
}

func TestGetByID(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)

	// テスト用のデータを追加
	dueDate := time.Now().Add(24 * time.Hour)
	todo := &domain.Todo{
		Title:   "テスト用TODO",
		DueDate: dueDate,
	}
	createdTodo, _ := repo.Create(todo)

	// テスト実行
	result, err := service.GetByID(createdTodo.ID)

	// アサーション
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	if result == nil {
		t.Error("Expected todo but got nil")
	}
	if result.ID != createdTodo.ID {
		t.Errorf("Expected ID %d but got %d", createdTodo.ID, result.ID)
	}
	if result.Title != "テスト用TODO" {
		t.Errorf("Expected title 'テスト用TODO' but got '%s'", result.Title)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)

	// 存在しないIDで検索
	_, err := service.GetByID(999)

	// アサーション
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestCreate(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)

	// テスト実行
	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	createdTodo, err := service.Create("新しいTODO", dueDate)

	// アサーション
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}
	if createdTodo == nil {
		t.Error("Expected todo but got nil")
	}
	if createdTodo.ID <= 0 {
		t.Errorf("Expected positive ID but got %d", createdTodo.ID)
	}
	if createdTodo.Title != "新しいTODO" {
		t.Errorf("Expected title '新しいTODO' but got '%s'", createdTodo.Title)
	}
}

func TestCreateWithEmptyTitle(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)

	// 空のタイトルでテスト実行
	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	_, err := service.Create("", dueDate)

	// アサーション
	if err == nil {
		t.Error("Expected error but got nil")
	}
	if !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("Expected ErrEmptyTitle but got %v", err)
	}
}

func TestCreateWithInvalidDueDate(t *testing.T) {
	// モックリポジトリの準備
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)

	// 不正な日付形式でテスト実行
	_, err := service.Create("テストTODO", "不正な日付形式")

	// アサーション
	if err == nil {
		t.Error("Expected error but got nil")
	}
	if !errors.Is(err, ErrInvalidDueDate) {
		t.Errorf("Expected ErrInvalidDueDate but got %v", err)
	}
}
