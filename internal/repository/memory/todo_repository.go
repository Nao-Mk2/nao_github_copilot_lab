package memory

import (
	"errors"
	"sync"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/domain"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

var (
	// ErrTodoNotFound はTODOアイテムが見つからない場合のエラーです
	ErrTodoNotFound = errors.New("todo not found")
)

// MemoryTodoRepository はメモリ上でTODOアイテムを管理するリポジトリの実装です
type MemoryTodoRepository struct {
	todos  map[int]*domain.Todo
	mutex  sync.RWMutex
	nextID int
}

// NewMemoryTodoRepository は新しいMemoryTodoRepositoryのインスタンスを生成します
func NewMemoryTodoRepository() repository.TodoRepository {
	return &MemoryTodoRepository{
		todos:  make(map[int]*domain.Todo),
		nextID: 1,
	}
}

// GetByID は指定されたIDのTODOアイテムを取得します
func (r *MemoryTodoRepository) GetByID(id int) (*domain.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

// Create は新しいTODOアイテムを作成します
func (r *MemoryTodoRepository) Create(todo *domain.Todo) (*domain.Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// IDの設定
	todo.ID = r.nextID
	r.nextID++

	// 新しいオブジェクトをマップに追加
	r.todos[todo.ID] = todo

	return todo, nil
}
