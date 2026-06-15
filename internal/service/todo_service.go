package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"todolist/internal/apperrors"
	"todolist/internal/model"
	"todolist/internal/repository"
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func parseDueDate(raw string) (time.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return time.Time{}, apperrors.ErrDueDateInvalid
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, apperrors.ErrDueDateInvalid
	}
	return t, nil
}

func parseTodoID(raw string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || id <= 0 {
		return 0, apperrors.ErrTodoIDNotPositive
	}
	return id, nil
}

func (s *TodoService) Create(ctx context.Context, req model.CreateTodoRequest) (*model.Todo, error) {
	if strings.TrimSpace(req.Task) == "" {
		return nil, apperrors.ErrTaskTitleEmpty
	}
	due, err := parseDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, strings.TrimSpace(req.Task), due)
}

func (s *TodoService) GetByID(ctx context.Context, idFromPath string) (*model.Todo, error) {
	id, err := parseTodoID(idFromPath)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TodoService) List(ctx context.Context, includeCompleted bool) ([]model.Todo, error) {
	return s.repo.List(ctx, includeCompleted)
}

func (s *TodoService) Update(ctx context.Context, idFromPath string, req model.UpdateTodoRequest) (*model.Todo, error) {
	id, err := parseTodoID(idFromPath)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Task) == "" {
		return nil, apperrors.ErrTaskTitleEmpty
	}
	due, err := parseDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, strings.TrimSpace(req.Task), due, req.Completed)
}

func (s *TodoService) Delete(ctx context.Context, idFromPath string) error {
	id, err := parseTodoID(idFromPath)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func PublicErrorMessage(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, apperrors.ErrTodoNotFound):
		return "todo not found"
	case errors.Is(err, apperrors.ErrTaskTitleEmpty):
		return "task cannot be empty"
	case errors.Is(err, apperrors.ErrDueDateInvalid):
		return "invalid due_date"
	case errors.Is(err, apperrors.ErrTodoIDNotPositive):
		return "invalid id"
	default:
		return "internal server error"
	}
}
