package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"todolist/internal/apperrors"
)

func TestDBRepository_Create_List_SortedByDueDate(t *testing.T) {
	repo := NewDBRepository()
	ctx := context.Background()

	laterDue, err := repo.Create(ctx, "later", time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	earlierDue, err := repo.Create(ctx, "earlier", time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if laterDue.ID == earlierDue.ID {
		t.Fatalf("ids must be unique: %d", laterDue.ID)
	}

	list, err := repo.List(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("want 2 got %d", len(list))
	}
	if list[0].Task != "earlier" || list[1].Task != "later" {
		t.Fatalf("unexpected order: %#v", list)
	}
}

func TestDBRepository_List_ExcludesCompletedByDefault(t *testing.T) {
	repo := NewDBRepository()
	ctx := context.Background()

	_, err := repo.Create(ctx, "open", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	completed, err := repo.Create(ctx, "done", time.Now().Add(48*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Update(ctx, completed.ID, completed.Task, completed.DueDate, true); err != nil {
		t.Fatal(err)
	}

	incompleteOnly, err := repo.List(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(incompleteOnly) != 1 || incompleteOnly[0].Task != "open" {
		t.Fatalf("want only open task, got %#v", incompleteOnly)
	}

	withCompleted, err := repo.List(ctx, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(withCompleted) != 2 {
		t.Fatalf("want 2 got %d", len(withCompleted))
	}
}

func TestDBRepository_ConcurrentCreates_UniqueIDs(t *testing.T) {
	repo := NewDBRepository()
	ctx := context.Background()
	const goroutineCount = 100
	var waitGroup sync.WaitGroup
	waitGroup.Add(goroutineCount)
	for index := 0; index < goroutineCount; index++ {
		go func(offset int) {
			defer waitGroup.Done()
			if _, err := repo.Create(ctx, "task", time.Now().Add(time.Duration(offset)*time.Second)); err != nil {
				t.Error(err)
			}
		}(index)
	}
	waitGroup.Wait()

	all, err := repo.List(ctx, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != goroutineCount {
		t.Fatalf("want %d got %d", goroutineCount, len(all))
	}
	seenIDs := make(map[int]struct{})
	for _, todo := range all {
		if _, duplicate := seenIDs[todo.ID]; duplicate {
			t.Fatalf("duplicate id %d", todo.ID)
		}
		seenIDs[todo.ID] = struct{}{}
	}
}

func TestDBRepository_Get_Update_Delete_NotFound(t *testing.T) {
	repo := NewDBRepository()
	ctx := context.Background()

	if _, err := repo.GetByID(ctx, 99); err != apperrors.ErrTodoNotFound {
		t.Fatalf("want not found got %v", err)
	}
	if err := repo.Delete(ctx, 99); err != apperrors.ErrTodoNotFound {
		t.Fatalf("want not found got %v", err)
	}

	created, err := repo.Create(ctx, "x", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Update(ctx, 999, "nope", time.Now(), false); err != apperrors.ErrTodoNotFound {
		t.Fatalf("want not found got %v", err)
	}

	updated, err := repo.Update(ctx, created.ID, "y", created.DueDate, true)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Task != "y" || !updated.Completed {
		t.Fatalf("unexpected update: %#v", updated)
	}
	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, created.ID); err != apperrors.ErrTodoNotFound {
		t.Fatalf("want not found got %v", err)
	}
}
