package repository

import (
	"errors"
	"sync"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/model"
)

var (
	// ErrTodoNotFound はTODOアイテムが見つからない場合のエラー
	ErrTodoNotFound = errors.New("TODOアイテムが見つかりません")
)

// TodoRepository はTODOアイテムのリポジトリのインターフェース
type TodoRepository interface {
	// Save はTODOアイテムを保存します
	Save(todo *model.Todo) (*model.Todo, error)

	// FindByID は指定されたIDのTODOアイテムを取得します
	FindByID(id int) (*model.Todo, error)
}

// InMemoryTodoRepository はTODOアイテムのインメモリリポジトリの実装
type InMemoryTodoRepository struct {
	mutex  sync.RWMutex
	todos  map[int]*model.Todo
	nextID int
}

// NewInMemoryTodoRepository は新しいInMemoryTodoRepositoryインスタンスを作成します
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos:  make(map[int]*model.Todo),
		nextID: 1,
	}
}

// Save はTODOアイテムをメモリに保存します
func (r *InMemoryTodoRepository) Save(todo *model.Todo) (*model.Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 新規作成の場合
	if todo.ID == 0 {
		todo.ID = r.nextID
		r.nextID++
	}

	r.todos[todo.ID] = todo
	return todo, nil
}

// FindByID は指定されたIDのTODOアイテムをメモリから取得します
func (r *InMemoryTodoRepository) FindByID(id int) (*model.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}
