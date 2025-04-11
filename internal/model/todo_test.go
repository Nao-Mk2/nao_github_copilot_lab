package model_test

import (
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
)

func TestNewTodo(t *testing.T) {
	t.Run("有効なTodoの作成", func(t *testing.T) {
		title := "テストタスク"
		dueDate := time.Now().Add(24 * time.Hour) // 明日

		todo, err := model.NewTodo(title, dueDate)

		if err != nil {
			t.Fatalf("エラーは期待されていません: %v", err)
		}

		if todo.Title != title {
			t.Errorf("タイトルが一致しません, 期待値: %s, 実際値: %s", title, todo.Title)
		}

		if !todo.DueDate.Equal(dueDate) {
			t.Errorf("期限が一致しません, 期待値: %s, 実際値: %s", dueDate, todo.DueDate)
		}
	})

	t.Run("空のタイトルでエラー", func(t *testing.T) {
		title := ""
		dueDate := time.Now().Add(24 * time.Hour) // 明日

		_, err := model.NewTodo(title, dueDate)

		if err != model.ErrEmptyTitle {
			t.Errorf("期待されるエラー: %v, 実際のエラー: %v", model.ErrEmptyTitle, err)
		}
	})

	t.Run("過去の期限でエラー", func(t *testing.T) {
		title := "テストタスク"
		dueDate := time.Now().Add(-24 * time.Hour) // 昨日

		_, err := model.NewTodo(title, dueDate)

		if err != model.ErrInvalidDueDate {
			t.Errorf("期待されるエラー: %v, 実際のエラー: %v", model.ErrInvalidDueDate, err)
		}
	})
}

func TestTodoValidate(t *testing.T) {
	t.Run("有効なTodoの検証", func(t *testing.T) {
		todo := &model.Todo{
			ID:      1,
			Title:   "テストタスク",
			DueDate: time.Now().Add(24 * time.Hour),
		}

		err := todo.Validate()

		if err != nil {
			t.Errorf("エラーは期待されていません: %v", err)
		}
	})

	t.Run("空のタイトルでエラー", func(t *testing.T) {
		todo := &model.Todo{
			ID:      1,
			Title:   "",
			DueDate: time.Now().Add(24 * time.Hour),
		}

		err := todo.Validate()

		if err != model.ErrEmptyTitle {
			t.Errorf("期待されるエラー: %v, 実際のエラー: %v", model.ErrEmptyTitle, err)
		}
	})

	t.Run("過去の期限でエラー", func(t *testing.T) {
		todo := &model.Todo{
			ID:      1,
			Title:   "テストタスク",
			DueDate: time.Now().Add(-24 * time.Hour),
		}

		err := todo.Validate()

		if err != model.ErrInvalidDueDate {
			t.Errorf("期待されるエラー: %v, 実際のエラー: %v", model.ErrInvalidDueDate, err)
		}
	})
}
