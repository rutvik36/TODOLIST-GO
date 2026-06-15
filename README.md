# To-Do List REST API (Go)

REST API for creating, listing, updating, and deleting to-dos. Matches [`instructions.md`](instructions.md).

## Features

- CRUD over JSON; validation for task text and RFC3339 `due_date`
- List sorted by due date; completed hidden unless `?include_completed=true`
- **handler → service → repository** with constructor wiring in `cmd/server`
- `DBRepository` (`sync.RWMutex` map), `slog` for errors, graceful shutdown
- Sample curl traces with real responses: [`docs/api_responses.md`](docs/api_responses.md)

## Layout

```text
cmd/server/main.go
internal/
  apperrors/errors.go    # shared sentinel errors
  handler/todos_http.go  # HTTP + JSON
  model/todo_models.go
  repository/db_repo.go       # TodoRepository interface + DBRepository (in-process map)
  service/todo_service.go
pkg/README.md
```

## Run

```bash
go test ./...
go run ./cmd/server          # :8080, or PORT=3000
```

Environment: `PORT` (optional, default `8080`; a leading `:` is added if missing).

## API (curl)

```bash
curl -sS -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"task":"Notes","due_date":"2026-06-20T18:00:00Z"}'

curl -sS http://localhost:8080/todos/1
curl -sS 'http://localhost:8080/todos'
curl -sS 'http://localhost:8080/todos?include_completed=true'

curl -sS -X PUT http://localhost:8080/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{"task":"Notes","due_date":"2026-06-22T18:00:00Z","completed":true}'

curl -i -X DELETE http://localhost:8080/todos/1
```

Errors: `{"error":"message"}` with **400** / **404** / **500** as in the spec.

## Database later

Implement `repository.TodoRepository` in a new file (e.g. SQL), construct it in `main`, and pass it to `service.NewTodoService`.
