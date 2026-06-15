# To-Do Application in Go

## Overview

Build a RESTful To-Do application in Go that supports creating, reading, updating, listing, and deleting to-do items.

The application should expose HTTP APIs and maintain a collection of tasks.

---

## Functional Requirements

### 1. Create a To-Do Item

Create a new task with task details and a due date.

#### Request

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

#### Response

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

Retrieve a specific to-do item by its ID.

#### Request

```http
GET /todos/{id}
```

#### Response

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

Retrieve all to-do items.

#### Requirements

* Sort tasks by due date (earliest first).
* Exclude completed tasks by default.
* Allow inclusion of completed tasks through a query parameter.

#### Requests

Exclude completed tasks (default):

```http
GET /todos
```

Include completed tasks:

```http
GET /todos?include_completed=true
```

#### Response

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

Update task details, due date, or completion status.

#### Request

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

#### Response

```json
{
  "id": 1,
  "task": "Prepare complete system design notes",
  "due_date": "2026-06-22T18:00:00Z",
  "completed": true
}
```

---

### 5. Delete a To-Do Item

Delete a task by its ID.

#### Request

```http
DELETE /todos/{id}
```

#### Response

```http
204 No Content
```

---

## Data Model

```go
type Todo struct {
    ID        int       `json:"id"`
    Task      string    `json:"task"`
    DueDate   time.Time `json:"due_date"`
    Completed bool      `json:"completed"`
}
```

---

## API Summary

| Method | Endpoint    | Description        |
| ------ | ----------- | ------------------ |
| POST   | /todos      | Create a new to-do |
| GET    | /todos/{id} | Get a to-do by ID  |
| GET    | /todos      | List to-dos        |
| PUT    | /todos/{id} | Update a to-do     |
| DELETE | /todos/{id} | Delete a to-do     |

---

## Non-Functional Requirements

### Validation

* Task text must not be empty.
* Due date must be a valid timestamp.
* ID must be unique.

### Error Handling

#### Invalid Request

```json
{
  "error": "task cannot be empty"
}
```

#### Resource Not Found

```json
{
  "error": "todo not found"
}
```

### HTTP Status Codes

| Status Code | Meaning               |
| ----------- | --------------------- |
| 200         | Success               |
| 201         | Created               |
| 204         | Deleted               |
| 400         | Bad Request           |
| 404         | Not Found             |
| 500         | Internal Server Error |

---

## Suggested Project Structure

```text
todo-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── pkg/
├── go.mod
└── README.md
```

---

## Bonus Enhancements

* Persistent storage using SQLite or PostgreSQL.
* Pagination for list API.
* Search by task name.
* Unit tests.
* Docker support.
* Swagger/OpenAPI documentation.
* Authentication and authorization.

---

## Expected Behavior

1. Users can create tasks with a due date.
2. Tasks can be marked as completed.
3. Completed tasks are excluded from list responses by default.
4. Tasks are always returned sorted by due date in ascending order.
5. CRUD operations follow REST principles and proper HTTP status codes.
