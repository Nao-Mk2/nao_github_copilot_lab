package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

// TodoHandler はTODOアイテムのハンドラー
type TodoHandler struct {
	repo repository.TodoRepository
}

// NewTodoHandler は新しいTodoHandlerインスタンスを作成します
func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{
		repo: repo,
	}
}

// todoRequest はTODOアイテム作成リクエストの構造体
type todoRequest struct {
	Title   string `json:"title"`
	DueDate string `json:"due_date"` // RFC3339形式
}

// GetTodo は指定されたIDのTODOアイテムを取得するハンドラー
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	// URLからIDを抽出 (/todo/{id}の形式)
	path := r.URL.Path
	segments := strings.Split(path, "/")

	if len(segments) < 3 {
		http.Error(w, "不正なURL形式です", http.StatusBadRequest)
		return
	}

	idStr := segments[len(segments)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "不正なIDです", http.StatusBadRequest)
		return
	}

	// リポジトリからTODOアイテムを取得
	todo, err := h.repo.FindByID(id)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			http.Error(w, "TODOアイテムが見つかりません", http.StatusNotFound)
		} else {
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
		}
		return
	}

	// レスポンスを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

// CreateTodo は新しいTODOアイテムを作成するハンドラー
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// リクエストボディをデコード
	var req todoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "不正なリクエスト形式です", http.StatusBadRequest)
		return
	}

	// 期限の形式を検証およびパース
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		http.Error(w, "不正な期限形式です。RFC3339形式で指定してください", http.StatusBadRequest)
		return
	}

	// 新しいTODOアイテムを作成
	todo, err := model.NewTodo(req.Title, dueDate)
	if err != nil {
		var statusCode int
		var message string

		switch err {
		case model.ErrEmptyTitle:
			statusCode = http.StatusBadRequest
			message = err.Error()
		case model.ErrInvalidDueDate:
			statusCode = http.StatusBadRequest
			message = err.Error()
		default:
			statusCode = http.StatusInternalServerError
			message = "内部サーバーエラー"
		}

		http.Error(w, message, statusCode)
		return
	}

	// リポジトリに保存
	savedTodo, err := h.repo.Save(todo)
	if err != nil {
		http.Error(w, "TODOアイテムの保存に失敗しました", http.StatusInternalServerError)
		return
	}

	// レスポンスを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedTodo)
}
