package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/service"
)

// TodoHandler はTODOアイテムのHTTPハンドラーです
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler は新しいTodoHandlerのインスタンスを生成します
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodoRequest はTODOアイテム作成リクエストのJSONを表す構造体です
type CreateTodoRequest struct {
	Title   string `json:"title"`
	DueDate string `json:"due_date"`
}

// ServeHTTP はHTTPリクエストを処理するハンドラー関数です
func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// URLパスからIDを抽出するための処理
	path := strings.TrimPrefix(r.URL.Path, "/todo")

	// Content-TypeをJSONに設定
	w.Header().Set("Content-Type", "application/json")

	if path == "" || path == "/" {
		// /todoへのリクエスト
		switch r.Method {
		case http.MethodPost:
			h.createTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /todo/{id}へのリクエスト
	id, err := strconv.Atoi(strings.TrimPrefix(path, "/"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTodo(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createTodo はTODOアイテムを作成するハンドラー関数です
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	todo, err := h.todoService.Create(req.Title, req.DueDate)
	if err != nil {
		// エラーによって適切なHTTPステータスを返す
		switch err {
		case service.ErrEmptyTitle:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Title cannot be empty"})
		case service.ErrInvalidDueDate:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid due date format"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
		return
	}

	// 成功レスポンス
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// getTodo は指定されたIDのTODOアイテムを取得するハンドラー関数です
func (h *TodoHandler) getTodo(w http.ResponseWriter, r *http.Request, id int) {
	todo, err := h.todoService.GetByID(id)
	if err != nil {
		// エラーによって適切なHTTPステータスを返す
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Todo with ID %d not found", id)})
		return
	}

	// 成功レスポンス
	json.NewEncoder(w).Encode(todo)
}
