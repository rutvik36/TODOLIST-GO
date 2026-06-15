package model

import "time"

// Todo is a single to-do item.
type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	DueDate   time.Time `json:"due_date"`
	Completed bool      `json:"completed"`
}

// CreateTodoRequest is the JSON body for POST /todos.
type CreateTodoRequest struct {
	Task    string `json:"task"`
	DueDate string `json:"due_date"`
}

// UpdateTodoRequest is the JSON body for PUT /todos/{id}.
type UpdateTodoRequest struct {
	Task      string `json:"task"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}
