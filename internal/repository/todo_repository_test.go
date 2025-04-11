package repository_test

import (
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

func TestInMemoryTodoRepository_Save(t *testing.T) {
	// リポジトリの作成
	repo := repository.NewInMemoryTodoRepository()

	// テストケース1: 新規Todoの保存（IDが自動で割り当てられる）
	dueDate := time.Now().Add(24 * time.Hour) // 明日
	todo, _ := model.NewTodo("テストタスク", dueDate)

	savedTodo, err := repo.Save(todo)
	if err != nil {
		t.Errorf("新規Todoの保存に失敗: %v", err)
	}
	if savedTodo.ID != 1 {
		t.Errorf("予期しないID。期待値: %d, 実際: %d", 1, savedTodo.ID)
	}
	if savedTodo.Title != "テストタスク" {
		t.Errorf("予期しないタイトル。期待値: %s, 実際: %s", "テストタスク", savedTodo.Title)
	}

	// テストケース2: 既存Todoの更新
	savedTodo.Title = "更新されたタスク"
	updatedTodo, err := repo.Save(savedTodo)
	if err != nil {
		t.Errorf("既存Todoの更新に失敗: %v", err)
	}
	if updatedTodo.ID != 1 {
		t.Errorf("予期しないID。期待値: %d, 実際: %d", 1, updatedTodo.ID)
	}
	if updatedTodo.Title != "更新されたタスク" {
		t.Errorf("予期しないタイトル。期待値: %s, 実際: %s", "更新されたタスク", updatedTodo.Title)
	}

	// テストケース3: 別のTodoの保存（IDが自動的に増加する）
	anotherTodo, _ := model.NewTodo("別のタスク", dueDate)
	anotherSavedTodo, err := repo.Save(anotherTodo)
	if err != nil {
		t.Errorf("別のTodoの保存に失敗: %v", err)
	}
	if anotherSavedTodo.ID != 2 {
		t.Errorf("予期しないID。期待値: %d, 実際: %d", 2, anotherSavedTodo.ID)
	}
}

func TestInMemoryTodoRepository_FindByID(t *testing.T) {
	// リポジトリの作成
	repo := repository.NewInMemoryTodoRepository()

	// テスト用のTodoを保存
	dueDate := time.Now().Add(24 * time.Hour) // 明日
	todo, _ := model.NewTodo("テストタスク", dueDate)
	savedTodo, _ := repo.Save(todo)

	// テストケース1: 存在するIDのTodoを取得
	foundTodo, err := repo.FindByID(savedTodo.ID)
	if err != nil {
		t.Errorf("存在するTodoの取得に失敗: %v", err)
	}
	if foundTodo.ID != savedTodo.ID {
		t.Errorf("予期しないID。期待値: %d, 実際: %d", savedTodo.ID, foundTodo.ID)
	}
	if foundTodo.Title != savedTodo.Title {
		t.Errorf("予期しないタイトル。期待値: %s, 実際: %s", savedTodo.Title, foundTodo.Title)
	}

	// テストケース2: 存在しないIDの場合はエラーを返す
	nonExistentID := 999
	_, err = repo.FindByID(nonExistentID)
	if err != repository.ErrTodoNotFound {
		t.Errorf("存在しないIDでエラーが返されなかった。期待値: %v, 実際: %v", repository.ErrTodoNotFound, err)
	}
}
