package memory

import (
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
)

func TestMemoryTodoRepository_Create(t *testing.T) {
	repo := NewMemoryTodoRepository()

	// テスト用のTodoアイテムを作成
	dueDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	todo := &domain.Todo{
		Title:   "テストタスク",
		DueDate: dueDate,
	}

	// リポジトリにTodoを作成
	createdTodo, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// 作成されたTodoの検証
	if createdTodo.ID <= 0 {
		t.Errorf("Create() generated ID = %v, want > 0", createdTodo.ID)
	}

	if createdTodo.Title != todo.Title {
		t.Errorf("Create() title = %v, want %v", createdTodo.Title, todo.Title)
	}

	if !createdTodo.DueDate.Equal(todo.DueDate) {
		t.Errorf("Create() dueDate = %v, want %v", createdTodo.DueDate, todo.DueDate)
	}
}

func TestMemoryTodoRepository_GetByID(t *testing.T) {
	repo := NewMemoryTodoRepository()

	// テスト用のTodoアイテムを作成して保存
	dueDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	todo := &domain.Todo{
		Title:   "テストタスク",
		DueDate: dueDate,
	}

	createdTodo, _ := repo.Create(todo)

	// 存在するIDのケース
	t.Run("存在するID", func(t *testing.T) {
		fetchedTodo, err := repo.GetByID(createdTodo.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if fetchedTodo.ID != createdTodo.ID {
			t.Errorf("GetByID() ID = %v, want %v", fetchedTodo.ID, createdTodo.ID)
		}

		if fetchedTodo.Title != createdTodo.Title {
			t.Errorf("GetByID() title = %v, want %v", fetchedTodo.Title, createdTodo.Title)
		}

		if !fetchedTodo.DueDate.Equal(createdTodo.DueDate) {
			t.Errorf("GetByID() dueDate = %v, want %v", fetchedTodo.DueDate, createdTodo.DueDate)
		}
	})

	// 存在しないIDのケース
	t.Run("存在しないID", func(t *testing.T) {
		nonExistentID := createdTodo.ID + 1000
		_, err := repo.GetByID(nonExistentID)
		if err != ErrTodoNotFound {
			t.Errorf("GetByID() error = %v, want %v", err, ErrTodoNotFound)
		}
	})
}
