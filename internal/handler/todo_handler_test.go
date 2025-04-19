package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/service"
)

// MockTodoService はテスト用のモックサービスです
type MockTodoService struct {
	GetByIDFunc func(id int) (*domain.Todo, error)
	CreateFunc  func(title string, dueDate string) (*domain.Todo, error)
}

func (m *MockTodoService) GetByID(id int) (*domain.Todo, error) {
	return m.GetByIDFunc(id)
}

func (m *MockTodoService) Create(title string, dueDate string) (*domain.Todo, error) {
	return m.CreateFunc(title, dueDate)
}

func TestGetTodo(t *testing.T) {
	// テスト用の日付
	dueDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	// テストケース
	testCases := []struct {
		name            string
		todoID          string
		mockGetByIDFunc func(id int) (*domain.Todo, error)
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:   "正常系 - TODOアイテムの取得",
			todoID: "1",
			mockGetByIDFunc: func(id int) (*domain.Todo, error) {
				return &domain.Todo{
					ID:      1,
					Title:   "テストTODO",
					DueDate: dueDate,
				}, nil
			},
			expectedStatus: http.StatusOK,
			// 期待されるJSONにはdueDate形式が反映されることに注意
			expectedBody: `{"id":1,"title":"テストTODO","due_date":"2025-12-31T23:59:59Z"}`,
		},
		{
			name:   "異常系 - TODOアイテムが存在しない",
			todoID: "999",
			mockGetByIDFunc: func(id int) (*domain.Todo, error) {
				return nil, errors.New("todo not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Todo with ID 999 not found"}`,
		},
		{
			name:            "異常系 - 不正なID",
			todoID:          "invalid",
			mockGetByIDFunc: nil, // この場合は呼ばれない
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `Invalid todo ID`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックサービスの準備
			mockService := &MockTodoService{
				GetByIDFunc: tc.mockGetByIDFunc,
			}

			// ハンドラーの生成
			handler := NewTodoHandler(mockService)

			// リクエストの作成
			req, err := http.NewRequest(http.MethodGet, "/todo/"+tc.todoID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスレコーダーの作成
			rr := httptest.NewRecorder()

			// リクエストの実行
			handler.ServeHTTP(rr, req)

			// ステータスコードの検証
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("期待したステータスコード %v が %v と一致しません", tc.expectedStatus, status)
			}

			// レスポンスボディの検証（BadRequestの場合は部分一致チェック）
			if tc.expectedStatus == http.StatusBadRequest && tc.todoID == "invalid" {
				if rr.Body.String() != "Invalid todo ID\n" {
					t.Errorf("期待したレスポンスボディ %q が %q と一致しません", "Invalid todo ID", rr.Body.String())
				}
			} else {
				// JSONレスポンスの場合は改行を削除して比較
				actual := rr.Body.String()
				actual = actual[:len(actual)-1] // 末尾の改行を削除
				if actual != tc.expectedBody {
					t.Errorf("期待したレスポンスボディ %q が %q と一致しません", tc.expectedBody, actual)
				}
			}
		})
	}
}

func TestCreateTodo(t *testing.T) {
	// テスト用の日付
	dueDate := "2025-12-31T23:59:59Z"
	parsedDueDate, _ := time.Parse(time.RFC3339, dueDate)

	// テストケース
	testCases := []struct {
		name           string
		requestBody    string
		mockCreateFunc func(title string, dueDate string) (*domain.Todo, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系 - TODOアイテムの作成",
			requestBody: `{"title":"新しいTODO","due_date":"2025-12-31T23:59:59Z"}`,
			mockCreateFunc: func(title string, dueDate string) (*domain.Todo, error) {
				return &domain.Todo{
					ID:      1,
					Title:   title,
					DueDate: parsedDueDate,
				}, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"title":"新しいTODO","due_date":"2025-12-31T23:59:59Z"}`,
		},
		{
			name:        "異常系 - タイトルが空",
			requestBody: `{"title":"","due_date":"2025-12-31T23:59:59Z"}`,
			mockCreateFunc: func(title string, dueDate string) (*domain.Todo, error) {
				return nil, service.ErrEmptyTitle
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Title cannot be empty"}`,
		},
		{
			name:        "異常系 - 不正な日付形式",
			requestBody: `{"title":"新しいTODO","due_date":"不正な日付"}`,
			mockCreateFunc: func(title string, dueDate string) (*domain.Todo, error) {
				return nil, service.ErrInvalidDueDate
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid due date format"}`,
		},
		{
			name:        "異常系 - 不正なリクエストボディ",
			requestBody: `不正なJSON`,
			mockCreateFunc: func(title string, dueDate string) (*domain.Todo, error) {
				return nil, nil // この場合は呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックサービスの準備
			mockService := &MockTodoService{
				CreateFunc: tc.mockCreateFunc,
			}

			// ハンドラーの生成
			handler := NewTodoHandler(mockService)

			// リクエストの作成
			req, err := http.NewRequest(http.MethodPost, "/todo", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			rr := httptest.NewRecorder()

			// リクエストの実行
			handler.ServeHTTP(rr, req)

			// ステータスコードの検証
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("期待したステータスコード %v が %v と一致しません", tc.expectedStatus, status)
			}

			// レスポンスボディの検証
			actual := rr.Body.String()
			actual = actual[:len(actual)-1] // 末尾の改行を削除
			if actual != tc.expectedBody {
				t.Errorf("期待したレスポンスボディ %q が %q と一致しません", tc.expectedBody, actual)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	// テストケース
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "異常系 - /todoにGETリクエスト",
			method:         http.MethodGet,
			path:           "/todo",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "異常系 - /todo/1にPOSTリクエスト",
			method:         http.MethodPost,
			path:           "/todo/1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "異常系 - /todoにPUTリクエスト",
			method:         http.MethodPut,
			path:           "/todo",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "異常系 - /todo/1にDELETEリクエスト",
			method:         http.MethodDelete,
			path:           "/todo/1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックサービスの準備
			mockService := &MockTodoService{}

			// ハンドラーの生成
			handler := NewTodoHandler(mockService)

			// リクエストの作成
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスレコーダーの作成
			rr := httptest.NewRecorder()

			// リクエストの実行
			handler.ServeHTTP(rr, req)

			// ステータスコードの検証
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("期待したステータスコード %v が %v と一致しません", tc.expectedStatus, status)
			}
		})
	}
}
