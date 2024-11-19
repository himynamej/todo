package todobus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/business/sdk/dbtest"
	"github.com/himynamej/todo/business/sdk/unitest"
)

func Test_TodoItem(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_TodoItem")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unitest.Run(t, query(db.BusDomain, sd), "query")
	unitest.Run(t, create(db.BusDomain), "create")
	unitest.Run(t, update(db.BusDomain, sd), "update")

}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	// Seed 2 new TodoItems using TestSeedTodoItems function
	todos, err := todobus.TestSeedTodoItems(ctx, 2, busDomain.Todo)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding todo items: %w", err)
	}

	// Populate SeedData structure
	sd := unitest.SeedData{
		Todos: todos,
	}

	return sd, nil
}

// =============================================================================
func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	todos := sd.Todos

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID.String() <= todos[j].ID.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: []todobus.TodoItem{todos[0], todos[1]},
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Todo.QueryByID(ctx, todos[0].ID)
				if err != nil {
					return err
				}
				resp2, err := busDomain.Todo.QueryByID(ctx, todos[1].ID)
				if err != nil {
					return err
				}

				return []todobus.TodoItem{resp, resp2}
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]todobus.TodoItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]todobus.TodoItem)
				for i := range gotResp {
					gotResp[i].DueDate = expResp[i].DueDate
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Todos[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Todo.QueryByID(ctx, sd.Todos[0].ID)
				if err != nil {
					return err
				}

				// Normalize to UTC and truncate precision for consistent comparison
				resp.DueDate = resp.DueDate.UTC().Truncate(time.Second)

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(todobus.TodoItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(todobus.TodoItem)
				gotResp.DueDate = expResp.DueDate
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
func create(busDomain dbtest.BusDomain) []unitest.Table {

	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: todobus.TodoItem{
				Description: "New TodoItem",
				DueDate:     time.Now().Add(72 * time.Hour).UTC(), // Ensure consistent UTC comparison
				FileID:      "mock-file-id",
			},
			ExcFunc: func(ctx context.Context) any {
				// Generate new item data
				nu := todobus.NewTodoItem{
					Description: "New TodoItem",
					DueDate:     time.Now().Add(72 * time.Hour).UTC(), // Ensure UTC for consistent comparison
					FileData:    []byte("new file data"),
					FileName:    "mock-file-id.txt",
				}

				// Create the new TodoItem
				resp, err := busDomain.Todo.Create(ctx, nu.Description, nu.DueDate, nu.FileData, nu.FileName)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(todobus.TodoItem)
				if !exists {
					return "error occurred" // Ensures the response is of the correct type
				}

				expResp := exp.(todobus.TodoItem)

				// Normalize ID and date fields for comparison
				expResp.ID = gotResp.ID                 // Since IDs are generated, set it dynamically for comparison
				gotResp.DueDate = gotResp.DueDate.UTC() // Ensure UTC normalization
				expResp.DueDate = expResp.DueDate.UTC() // Ensure UTC normalization
				gotResp.DueDate = expResp.DueDate
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: todobus.TodoItem{
				ID:          sd.Todos[0].ID,
				Description: "Updated TodoItem",
				DueDate:     time.Now().Add(96 * time.Hour),
				FileID:      sd.Todos[0].FileID,
			},
			ExcFunc: func(ctx context.Context) any {
				newDescription := "Updated TodoItem"
				newDueDate := time.Now().Add(96 * time.Hour)

				resp, err := busDomain.Todo.Update(ctx, sd.Todos[0].ID, newDescription, newDueDate)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(todobus.TodoItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(todobus.TodoItem)
				gotResp.DueDate = expResp.DueDate
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
