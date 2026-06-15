package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"todolist/internal/apperrors"
	"todolist/internal/model"
)

// TodoRepository persists todos. Implementations must be safe for concurrent use.
type TodoRepository interface {
	Create(ctx context.Context, task string, dueDate time.Time) (*model.Todo, error)
	GetByID(ctx context.Context, id int) (*model.Todo, error)
	List(ctx context.Context, includeCompleted bool) ([]model.Todo, error)
	Update(ctx context.Context, id int, task string, dueDate time.Time, completed bool) (*model.Todo, error)
	Delete(ctx context.Context, id int) error
}

// DBRepository implements TodoRepository with a mutex-backed map (in-process; replace with SQL when needed).
type DBRepository struct {
	mu         sync.RWMutex
	byID       map[int]*model.Todo
	nextTodoID int
}

func NewDBRepository() *DBRepository {
	return &DBRepository{
		byID:       make(map[int]*model.Todo),
		nextTodoID: 1,
	}
}

func (r *DBRepository) Create(_ context.Context, task string, dueDate time.Time) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextTodoID
	r.nextTodoID++

	t := &model.Todo{ID: id, Task: task, DueDate: dueDate, Completed: false}
	r.byID[id] = t
	out := *t
	return &out, nil
}

func (r *DBRepository) GetByID(_ context.Context, id int) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.byID[id]
	if !ok {
		return nil, apperrors.ErrTodoNotFound
	}
	out := *t
	return &out, nil
}

func (r *DBRepository) List(_ context.Context, includeCompleted bool) ([]model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []model.Todo
	for _, t := range r.byID {
		if !includeCompleted && t.Completed {
			continue
		}
		item := *t
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DueDate.Equal(out[j].DueDate) {
			return out[i].ID < out[j].ID
		}
		return out[i].DueDate.Before(out[j].DueDate)
	})
	return out, nil
}

func (r *DBRepository) Update(_ context.Context, id int, task string, dueDate time.Time, completed bool) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.byID[id]
	if !ok {
		return nil, apperrors.ErrTodoNotFound
	}
	t.Task = task
	t.DueDate = dueDate
	t.Completed = completed
	out := *t
	return &out, nil
}

func (r *DBRepository) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return apperrors.ErrTodoNotFound
	}
	delete(r.byID, id)
	return nil
}
