package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/handler"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

// テスト用のtodoRequest構造体 (オリジナルパッケージのものと同等)
type todoRequest struct {
	Title   string `json:"title"`
	DueDate string `json:"due_date"` // RFC3339形式
}

// 疑似リポジトリの実装（テスト用）
type mockTodoRepository struct {
	todos  map[int]*model.Todo
	nextID int
}

func newMockTodoRepository() *mockTodoRepository {
	return &mockTodoRepository{
		todos:  make(map[int]*model.Todo),
		nextID: 1,
	}
}

func (r *mockTodoRepository) Save(todo *model.Todo) (*model.Todo, error) {
	if todo.ID == 0 {
		todo.ID = r.nextID
		r.nextID++
	}
	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *mockTodoRepository) FindByID(id int) (*model.Todo, error) {
	todo, exists := r.todos[id]
	if !exists {
		return nil, repository.ErrTodoNotFound
	}
	return todo, nil
}

func TestGetTodo(t *testing.T) {
	// テスト用のリポジトリとハンドラーを準備
	mockRepo := newMockTodoRepository()
	todoHandler := handler.NewTodoHandler(mockRepo)

	// テスト用のTODOアイテムを作成してリポジトリに保存
	dueDate := time.Now().Add(24 * time.Hour)
	todo, _ := model.NewTodo("テストタスク", dueDate)
	savedTodo, _ := mockRepo.Save(todo)

	// テスト用のHTTPリクエストを作成
	req, _ := http.NewRequest("GET", "/todo/"+strconv.Itoa(savedTodo.ID), nil)
	responseRecorder := httptest.NewRecorder()

	// ハンドラーを実行
	todoHandler.GetTodo(responseRecorder, req)

	// レスポンスの検証
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("期待したステータスコードは %d、実際のステータスコードは %d", http.StatusOK, responseRecorder.Code)
	}

	var responseTodo model.Todo
	json.Unmarshal(responseRecorder.Body.Bytes(), &responseTodo)

	if responseTodo.ID != savedTodo.ID {
		t.Errorf("期待したIDは %d、実際のIDは %d", savedTodo.ID, responseTodo.ID)
	}

	if responseTodo.Title != savedTodo.Title {
		t.Errorf("期待したタイトルは %s、実際のタイトルは %s", savedTodo.Title, responseTodo.Title)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	// テスト用のリポジトリとハンドラーを準備
	mockRepo := newMockTodoRepository()
	todoHandler := handler.NewTodoHandler(mockRepo)

	// 存在しないIDでリクエスト
	req, _ := http.NewRequest("GET", "/todo/999", nil)
	responseRecorder := httptest.NewRecorder()

	// ハンドラーを実行
	todoHandler.GetTodo(responseRecorder, req)

	// レスポンスの検証
	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("期待したステータスコードは %d、実際のステータスコードは %d", http.StatusNotFound, responseRecorder.Code)
	}
}

func TestCreateTodo(t *testing.T) {
	// テスト用のリポジトリとハンドラーを準備
	mockRepo := newMockTodoRepository()
	todoHandler := handler.NewTodoHandler(mockRepo)

	// テスト用のリクエストボディを作成
	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	reqBody := todoRequest{
		Title:   "テスト作成タスク",
		DueDate: dueDate,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(body))
	responseRecorder := httptest.NewRecorder()

	// ハンドラーを実行
	todoHandler.CreateTodo(responseRecorder, req)

	// レスポンスの検証
	if responseRecorder.Code != http.StatusCreated {
		t.Errorf("期待したステータスコードは %d、実際のステータスコードは %d", http.StatusCreated, responseRecorder.Code)
	}

	var responseTodo model.Todo
	json.Unmarshal(responseRecorder.Body.Bytes(), &responseTodo)

	if responseTodo.ID != 1 {
		t.Errorf("期待したIDは %d、実際のIDは %d", 1, responseTodo.ID)
	}

	if responseTodo.Title != reqBody.Title {
		t.Errorf("期待したタイトルは %s、実際のタイトルは %s", reqBody.Title, responseTodo.Title)
	}
}

func TestCreateTodoInvalidRequest(t *testing.T) {
	// テスト用のリポジトリとハンドラーを準備
	mockRepo := newMockTodoRepository()
	todoHandler := handler.NewTodoHandler(mockRepo)

	// 無効な期限形式でリクエストボディを作成
	reqBody := todoRequest{
		Title:   "テスト作成タスク",
		DueDate: "無効な日付形式",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(body))
	responseRecorder := httptest.NewRecorder()

	// ハンドラーを実行
	todoHandler.CreateTodo(responseRecorder, req)

	// レスポンスの検証
	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("期待したステータスコードは %d、実際のステータスコードは %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestCreateTodoEmptyTitle(t *testing.T) {
	// テスト用のリポジトリとハンドラーを準備
	mockRepo := newMockTodoRepository()
	todoHandler := handler.NewTodoHandler(mockRepo)

	// 空のタイトルでリクエストボディを作成
	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	reqBody := todoRequest{
		Title:   "",
		DueDate: dueDate,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(body))
	responseRecorder := httptest.NewRecorder()

	// ハンドラーを実行
	todoHandler.CreateTodo(responseRecorder, req)

	// レスポンスの検証
	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("期待したステータスコードは %d、実際のステータスコードは %d", http.StatusBadRequest, responseRecorder.Code)
	}
}
