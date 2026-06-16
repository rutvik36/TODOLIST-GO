# To-Do Application in Go

This document describes the **product requirements** for the REST To-Do API and how **this repository** implements them. The code lives at **https://github.com/rutvik36/TODOLIST-GO** (Go module: `todolist`).

---

## Overview

A RESTful To-Do service in Go: create, read, list, update, and delete tasks over HTTP/JSON. Tasks are held in memory for the current implementation (see **Persistence** below).

---

## Functional Requirements

### 1. Create a To-Do Item

Create a new task with task text and a due date.

**Request**

```http
POST /todos
Content-Type: application/json
```

```json
{
  "task": "Prepare system design interview notes",
  "due_date": "2026-06-20T18:00:00Z"
}
```

**Response:** `201 Created`, JSON body with assigned `id` and `completed: false`.

```json
{
  "id": 1,
  "task": "Prepare system design interview notes",
  "due_date": "2026-06-20T18:00:00Z",
  "completed": false
}
```

---

### 2. Read a To-Do Item

Retrieve one item by numeric ID in the path.

**Request**

```http
GET /todos/{id}
```

**Response:** `200 OK` and JSON body, or `404 Not Found` with `{"error":"todo not found"}`.

```json
{
  "id": 1,
  "task": "Prepare system design interview notes",
  "due_date": "2026-06-20T18:00:00Z",
  "completed": false
}
```

---

### 3. List To-Do Items

**Behaviour**

* Sort by **due date ascending** (earliest first). If two items share the same due date, sort by **`id` ascending**.
* **Exclude** completed tasks by default.
* **`include_completed` query parameter:** when `true`, include completed items (same sort order).

**Requests**

```http
GET /todos
GET /todos?include_completed=true
```

**Response:** `200 OK`, JSON array (empty `[]` when there are no matches).

```json
[
  {
    "id": 2,
    "task": "Review Golang channels",
    "due_date": "2026-06-18T10:00:00Z",
    "completed": false
  },
  {
    "id": 1,
    "task": "Prepare system design interview notes",
    "due_date": "2026-06-20T18:00:00Z",
    "completed": false
  }
]
```

---

### 4. Update a To-Do Item

Replace task text, due date, and completion flag.

**Request**

```http
PUT /todos/{id}
Content-Type: application/json
```

```json
{
  "task": "Prepare complete system design notes",
  "due_date": "2026-06-22T18:00:00Z",
  "completed": true
}
```

**Response:** `200 OK` with full updated entity, or `404` / `400` as below.

---

### 5. Delete a To-Do Item

**Request**

```http
DELETE /todos/{id}
```

**Response:** `204 No Content` on success, or `404` with `{"error":"todo not found"}`.

---

## Data Model

Implemented in `internal/model/todo_models.go`:

```go
type Todo struct {
    ID        int       `json:"id"`
    Task      string    `json:"task"`
    DueDate   time.Time `json:"due_date"`
    Completed bool      `json:"completed"`
}
```

Request DTOs: `CreateTodoRequest` (`task`, `due_date` strings), `UpdateTodoRequest` (`task`, `due_date`, `completed`).

---

## API Summary

| Method | Endpoint    | Description        | Typical status |
| ------ | ----------- | ------------------ | -------------- |
| POST   | /todos      | Create a new to-do | 201            |
| GET    | /todos/{id} | Get a to-do by ID  | 200, 404       |
| GET    | /todos      | List to-dos        | 200            |
| PUT    | /todos/{id} | Update a to-do     | 200, 400, 404  |
| DELETE | /todos/{id} | Delete a to-do     | 204, 404       |

---

## Non-Functional Requirements

### Validation (implemented)

| Rule | Behaviour |
| ---- | --------- |
| Task text | Non-empty after trim; otherwise **400** `task cannot be empty`. |
| Due date | Must parse as **RFC3339** (e.g. `2026-06-20T18:00:00Z`); otherwise **400** `invalid due_date`. |
| Path `id` | Must be a positive integer; otherwise **400** `invalid id`. |
| Unique `id` | Assigned by the repository (monotonic); guaranteed unique while the process runs. |
| JSON body | Max **1 MiB**; unknown JSON fields rejected; malformed body → **400** `invalid request body`. |
| `include_completed` | If present and not a valid boolean → **400** `invalid include_completed`. |

### Error handling

Errors are JSON: `{ "error": "<message>" }`.

| HTTP | Example `error` values |
| ---- | ---------------------- |
| 400  | `task cannot be empty`, `invalid due_date`, `invalid id`, `invalid request body`, `invalid include_completed` |
| 404  | `todo not found` |
| 500  | `internal server error` (details only in server logs) |

### HTTP status codes

| Status | Meaning               |
| ------ | --------------------- |
| 200    | Success               |
| 201    | Created               |
| 204    | Deleted               |
| 400    | Bad Request           |
| 404    | Not Found             |
| 500    | Internal Server Error |

---

## Implementation (this repository)

### Architecture

* **`cmd/server`** — Wires `repository.DBRepository` → `service.TodoService` → `handler.TodoHTTPHandler`, starts `http.Server`, **graceful shutdown** on `SIGINT` / `SIGTERM` (10s shutdown timeout).
* **`internal/handler`** — HTTP only: routing (Go 1.22+ `http.ServeMux` method patterns), JSON encode/decode, status codes, delegates to service.
* **`internal/service`** — Validation, path ID parsing, `PublicErrorMessage` for stable client strings.
* **`internal/repository`** — `TodoRepository` interface and **`DBRepository`**: in-process `map` + `sync.RWMutex` (not a real SQL database; name reflects the persistence boundary).
* **`internal/apperrors`** — Sentinel errors shared by service and repository.
* **`internal/model`** — `Todo` and request structs.

### Project layout

```text
.
├── cmd/server/main.go
├── docs/api_responses.md     # captured curl -i traces
├── go.mod                    # module todolist, Go 1.22+
├── internal/
│   ├── apperrors/errors.go
│   ├── handler/todos_http.go
│   ├── handler/todos_http_test.go
│   ├── model/todo_models.go
│   ├── repository/db_repo.go        # TodoRepository + DBRepository
│   ├── repository/db_repo_test.go
│   ├── service/
│   │   ├── todo_service.go
│   │   └── todo_service_test.go
├── pkg/README.md
├── README.md
└── instructions.md           # this file
```

### Run and test

```bash
go test ./...
go run ./cmd/server
```

Listen address: environment variable **`PORT`** (default **`8080`**). If `PORT` does not start with `:`, a leading `:` is prepended (e.g. `8080` → `:8080`).

### Logging

Structured logs via **`log/slog`** (JSON to stdout): listen address, server/shutdown errors, and handler errors for **500** paths.

### Persistence

Data is **in memory only**; restarting the process clears all todos. To add a real database, implement `repository.TodoRepository` in a new file, construct it in `cmd/server/main.go`, and pass it to `service.NewTodoService`.

### Evidence of HTTP behaviour

See **`docs/api_responses.md`** for real `curl -i` output (status line, headers, bodies) against a local run.

---

## Bonus Enhancements (relative to this repo)

| Item | Status |
| ---- | ------ |
| Unit tests | **Done** — `repository`, `service`, `handler` (httptest). |
| Persistent storage (SQLite/PostgreSQL) | Not implemented; extension point is `TodoRepository`. |
| Pagination / search | Not implemented. |
| Docker / Swagger / auth | Not implemented. |

---

## Expected Behaviour (checklist)

1. Users can create tasks with a due date (**RFC3339**).
2. Tasks can be marked completed via **PUT**.
3. Completed tasks are omitted from **GET /todos** unless **`include_completed=true`**.
4. List responses are sorted by **due date ascending** (then **`id` ascending**).
5. CRUD uses the HTTP methods and status codes in this document.
