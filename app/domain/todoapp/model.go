package todoapp

import (
	"encoding/json"
	"time"

	"github.com/himynamej/todo/app/sdk/errs"
	"github.com/himynamej/todo/business/domain/todobus"
)

type queryParams struct {
	ID string
}

// FileUploadResponse represents the response returned when a file is uploaded.
type FileUploadResponse struct {
	FileID string `json:"fileId"`
}

// Encode implements the encoder interface for FileUploadResponse.
func (resp FileUploadResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(resp)
	return data, "application/json", err
}

// FileResponse represents the response returned for a file download.
type FileResponse struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

// Encode implements the encoder interface for FileResponse.
func (fileResp FileResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(fileResp)
	return data, "application/json", err
}

// TodoItem represents the structure for a Todo item in the application layer.
type TodoItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	FileID      string `json:"fileId"`
}

// Encode implements the encoder interface for a TodoItem.
func (app TodoItem) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// NewTodoItem defines the data needed to create a new TodoItem.
type NewTodoItem struct {
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"dueDate" validate:"required"`
	FileID      string `json:"fileId"`
}

// Encode implements the encoder interface.
func (app NewTodoItem) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// Decode implements the decoder interface.
func (app *NewTodoItem) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewTodoItem) Validate() error {
	if err := errs.Check(app); err != nil {
		return errs.Newf(errs.InvalidArgument, "validate: %s", err)
	}

	return nil
}

func toAppTodoItem(bus todobus.TodoItem) TodoItem {
	return TodoItem{
		ID:          bus.ID.String(),
		Description: bus.Description,
		DueDate:     bus.DueDate.Format(time.RFC3339),
		FileID:      bus.FileID,
	}
}
