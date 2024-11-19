package itemdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/himynamej/todo/business/domain/todobus"
)

// dbTodoItem represents the database structure of a TodoItem.
type dbTodoItem struct {
	ID          string    `db:"item_id"`
	Description string    `db:"description"`
	DueDate     time.Time `db:"due_date"`
	FileID      string    `db:"file_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

// toDBTodoItem converts a business-level TodoItem to a database-level TodoItem.
func toDBTodoItem(item todobus.TodoItem) dbTodoItem {
	return dbTodoItem{
		ID:          item.ID.String(),
		Description: item.Description,
		DueDate:     item.DueDate,
		FileID:      item.FileID,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
}

// toBusTodoItem converts a database-level TodoItem to a business-level TodoItem.
func toBusTodoItem(dbItem dbTodoItem) (todobus.TodoItem, error) {
	id, err := uuid.Parse(dbItem.ID)
	if err != nil {
		return todobus.TodoItem{}, fmt.Errorf("parse UUID: %w", err)
	}

	return todobus.TodoItem{
		ID:          id,
		Description: dbItem.Description,
		DueDate:     dbItem.DueDate,
		FileID:      dbItem.FileID,
	}, nil
}

// toBusTodoItems converts a slice of database-level TodoItems to a slice of business-level TodoItems.
func toBusTodoItems(dbItems []dbTodoItem) ([]todobus.TodoItem, error) {
	items := make([]todobus.TodoItem, len(dbItems))
	for i, dbItem := range dbItems {
		busItem, err := toBusTodoItem(dbItem)
		if err != nil {
			return nil, err
		}
		items[i] = busItem
	}
	return items, nil
}
