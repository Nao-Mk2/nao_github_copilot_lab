package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/handler"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

// テスト用のサーバーセットアップ関数
func setupTestServer() (*httptest.Server, *repository.InMemoryTodoRepository) {
	// リポジトリの初期化
	todoRepo := repository.NewInMemoryTodoRepository()

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoRepo)

	// ルーティングの設定
	mux := http.NewServeMux()

	// /todoエンドポイント（POST用）
	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			todoHandler.CreateTodo(w, r)
		} else {
			http.Error(w, "不正なメソッドです", http.StatusMethodNotAllowed)
		}
	})

	// /todo/パターンを含むエンドポイント（GET用）
	mux.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			todoHandler.GetTodo(w, r)
		} else {
			http.Error(w, "不正なメソッドまたはURLです", http.StatusMethodNotAllowed)
		}
	})

	// テスト用サーバーの起動
	server := httptest.NewServer(mux)
	return server, todoRepo
}

// TestCreateTodo はTODOアイテム作成APIのテスト
func TestCreateTodo(t *testing.T) {
	// テストサーバーのセットアップ
	server, _ := setupTestServer()
	defer server.Close()

	// 現在時刻から1日後の期限を設定
	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	// テスト用リクエストボディを作成
	reqBody := map[string]string{
		"title":    "テスト用TODOタスク",
		"due_date": dueDate,
	}
	reqBytes, _ := json.Marshal(reqBody)

	// POSTリクエストを送信
	resp, err := http.Post(
		fmt.Sprintf("%s/todo", server.URL),
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		t.Fatalf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを検証
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("期待するステータスコードは %d ですが、%d が返されました", http.StatusCreated, resp.StatusCode)
	}

	// レスポンスボディをデコード
	var todo model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatalf("レスポンスのデコードに失敗しました: %v", err)
	}

	// レスポンスの内容を検証
	if todo.ID <= 0 {
		t.Errorf("有効なIDが返されませんでした: %d", todo.ID)
	}
	if todo.Title != reqBody["title"] {
		t.Errorf("期待するタイトルは %s ですが、%s が返されました", reqBody["title"], todo.Title)
	}
}

// TestGetTodo はTODOアイテム取得APIのテスト
func TestGetTodo(t *testing.T) {
	// テストサーバーのセットアップ
	server, repo := setupTestServer()
	defer server.Close()

	// テスト用のTODOアイテムを作成
	dueDate := time.Now().Add(24 * time.Hour)
	todo, _ := model.NewTodo("サンプルTODO", dueDate)
	savedTodo, _ := repo.Save(todo)
	todoID := savedTodo.ID

	// GETリクエストを送信
	resp, err := http.Get(fmt.Sprintf("%s/todo/%d", server.URL, todoID))
	if err != nil {
		t.Fatalf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを検証
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待するステータスコードは %d ですが、%d が返されました", http.StatusOK, resp.StatusCode)
	}

	// レスポンスボディをデコード
	var responseTodo model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&responseTodo); err != nil {
		t.Fatalf("レスポンスのデコードに失敗しました: %v", err)
	}

	// レスポンスの内容を検証
	if responseTodo.ID != todoID {
		t.Errorf("期待するIDは %d ですが、%d が返されました", todoID, responseTodo.ID)
	}
	if responseTodo.Title != todo.Title {
		t.Errorf("期待するタイトルは %s ですが、%s が返されました", todo.Title, responseTodo.Title)
	}
}

// TestGetNonExistentTodo は存在しないTODOアイテムの取得テスト
func TestGetNonExistentTodo(t *testing.T) {
	// テストサーバーのセットアップ
	server, _ := setupTestServer()
	defer server.Close()

	// 存在しないIDでGETリクエストを送信
	nonExistentID := 9999
	resp, err := http.Get(fmt.Sprintf("%s/todo/%d", server.URL, nonExistentID))
	if err != nil {
		t.Fatalf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを検証
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期待するステータスコードは %d ですが、%d が返されました", http.StatusNotFound, resp.StatusCode)
	}
}

// TestCreateTodoWithInvalidData は不正なデータによるTODOアイテム作成のテスト
func TestCreateTodoWithInvalidData(t *testing.T) {
	// テストサーバーのセットアップ
	server, _ := setupTestServer()
	defer server.Close()

	testCases := []struct {
		name           string
		reqBody        map[string]string
		expectedStatus int
	}{
		{
			name: "空のタイトル",
			reqBody: map[string]string{
				"title":    "",
				"due_date": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "過去の期限",
			reqBody: map[string]string{
				"title":    "有効なタイトル",
				"due_date": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "不正な日付形式",
			reqBody: map[string]string{
				"title":    "有効なタイトル",
				"due_date": "2025-01-01", // RFC3339形式ではない
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBytes, _ := json.Marshal(tc.reqBody)

			// POSTリクエストを送信
			resp, err := http.Post(
				fmt.Sprintf("%s/todo", server.URL),
				"application/json",
				bytes.NewBuffer(reqBytes),
			)
			if err != nil {
				t.Fatalf("リクエストの送信に失敗しました: %v", err)
			}
			defer resp.Body.Close()

			// レスポンスのステータスコードを検証
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("期待するステータスコードは %d ですが、%d が返されました", tc.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// TestInvalidMethod は不正なHTTPメソッドのテスト
func TestInvalidMethod(t *testing.T) {
	// テストサーバーのセットアップ
	server, _ := setupTestServer()
	defer server.Close()

	// PUTリクエストを作成（APIでサポートされていないメソッド）
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/todo", server.URL), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを検証
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("期待するステータスコードは %d ですが、%d が返されました", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
