package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"todolist/internal/apperrors"
	"todolist/internal/model"
	"todolist/internal/repository"
)

func TestCreateValidation(t *testing.T) {
	s := NewTodoService(repository.NewDBRepository())
	ctx := context.Background()

	_, err := s.Create(ctx, model.CreateTodoRequest{Task: "", DueDate: "2026-06-20T18:00:00Z"})
	if !errors.Is(err, apperrors.ErrTaskTitleEmpty) {
		t.Fatalf("got %v", err)
	}
	_, err = s.Create(ctx, model.CreateTodoRequest{Task: "ok", DueDate: "bad"})
	if !errors.Is(err, apperrors.ErrDueDateInvalid) {
		t.Fatalf("got %v", err)
	}
	todo, err := s.Create(ctx, model.CreateTodoRequest{Task: "  x  ", DueDate: "2026-06-20T18:00:00Z"})
	if err != nil || todo.Task != "x" {
		t.Fatalf("%v %q", err, todo.Task)
	}
}

func TestGetByIDInvalid(t *testing.T) {
	s := NewTodoService(repository.NewDBRepository())
	if _, err := s.GetByID(context.Background(), "0"); !errors.Is(err, apperrors.ErrTodoIDNotPositive) {
		t.Fatalf("got %v", err)
	}
}

func TestListHidesCompleted(t *testing.T) {
	repo := repository.NewDBRepository()
	ctx := context.Background()
	_, _ = repo.Create(ctx, "open", time.Now().Add(time.Hour))
	done, _ := repo.Create(ctx, "done", time.Now().Add(2*time.Hour))
	_, _ = repo.Update(ctx, done.ID, done.Task, done.DueDate, true)

	s := NewTodoService(repo)
	list, err := s.List(ctx, false)
	if err != nil || len(list) != 1 {
		t.Fatalf("%v len=%d", err, len(list))
	}
}

func TestPublicErrorMessage(t *testing.T) {
	if PublicErrorMessage(apperrors.ErrTodoNotFound) != "todo not found" {
		t.Fatal()
	}
	if PublicErrorMessage(errors.New("x")) != "internal server error" {
		t.Fatal()
	}
}
