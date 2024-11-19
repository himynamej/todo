package todobus

import (
	"time"

	"github.com/google/uuid"
)

// TodoItem represents the structure for a todo item.
type TodoItem struct {
	ID          uuid.UUID
	Description string
	DueDate     time.Time
	FileID      string
}

// NewTodoItem contains information needed to create a new TodoItem.
type NewTodoItem struct {
	Description string
	DueDate     time.Time
	FileData    []byte
	FileName    string
}

// UpdateTodoItem contains information needed to update an existing TodoItem.
type UpdateTodoItem struct {
	Description *string
	DueDate     *time.Time
}
