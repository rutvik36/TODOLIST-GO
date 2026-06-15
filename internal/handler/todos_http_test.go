package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todolist/internal/repository"
	"todolist/internal/service"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	repo := repository.NewDBRepository()
	svc := service.NewTodoService(repo)
	h := NewTodoHTTPHandler(svc, slog.New(slog.NewTextHandler(nopWriter{}, nil)))
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestCRUD(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	c := ts.Client()

	res, err := c.Post(ts.URL+"/todos", "application/json", bytes.NewBufferString(
		`{"task":"Prepare notes","due_date":"2026-06-20T18:00:00Z"}`,
	))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status %d", res.StatusCode)
	}
	var created map[string]any
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	id := int(created["id"].(float64))

	res, err = c.Get(ts.URL + "/todos/" + strconv.Itoa(id))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get %d", res.StatusCode)
	}

	body := `{"task":"Prepare notes","due_date":"2026-06-20T18:00:00Z","completed":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, ts.URL+"/todos/"+strconv.Itoa(id), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	res, err = c.Get(ts.URL + "/todos")
	if err != nil {
		t.Fatal(err)
	}
	var list []map[string]any
	json.NewDecoder(res.Body).Decode(&list)
	res.Body.Close()
	if len(list) != 0 {
		t.Fatalf("list len %d", len(list))
	}

	res, err = c.Get(ts.URL + "/todos?include_completed=true")
	if err != nil {
		t.Fatal(err)
	}
	json.NewDecoder(res.Body).Decode(&list)
	res.Body.Close()
	if len(list) != 1 {
		t.Fatalf("with completed len %d", len(list))
	}

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodDelete, ts.URL+"/todos/"+strconv.Itoa(id), nil)
	res, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("delete %d", res.StatusCode)
	}

	res, err = c.Get(ts.URL + "/todos/" + strconv.Itoa(id))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("after delete %d", res.StatusCode)
	}
}
