package todobus

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// TestNewTodoItems is a helper method for generating new todo items with random data for testing.
func TestNewTodoItems(n int) []NewTodoItem {
	newItems := make([]NewTodoItem, n)

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++

		item := NewTodoItem{
			Description: fmt.Sprintf("Description%d", idx),
			DueDate:     time.Now().Add(time.Duration(rand.Intn(100)) * time.Hour),
			FileData:    []byte(fmt.Sprintf("Test file data %d", idx)),
			FileName:    fmt.Sprintf("TestFile%d.txt", idx),
		}

		newItems[i] = item
	}

	return newItems
}

// TestSeedTodoItems is a helper method for seeding TodoItem data into the system for testing.
func TestSeedTodoItems(ctx context.Context, n int, api *Business) ([]TodoItem, error) {
	newItems := TestNewTodoItems(n)

	items := make([]TodoItem, len(newItems))
	for i, newItem := range newItems {
		item, err := api.Create(ctx, newItem.Description, newItem.DueDate, newItem.FileData, newItem.FileName)
		if err != nil {
			return nil, fmt.Errorf("seeding todo item: idx: %d : %w", i, err)
		}

		items[i] = item
	}

	return items, nil
}
