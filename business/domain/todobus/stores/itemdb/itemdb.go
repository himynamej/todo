// Package itemdb contains TodoItem related CRUD functionality.
package itemdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/business/sdk/sqldb"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for TodoItem database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (todobus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create inserts a new TodoItem into the database.
func (s *Store) Create(ctx context.Context, item todobus.TodoItem) error {
	const q = `
	INSERT INTO todo_items
		(item_id, description, due_date, file_id, date_created, date_updated)
	VALUES
		(:item_id, :description, :due_date, :file_id, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTodoItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies an existing TodoItem in the database.
func (s *Store) Update(ctx context.Context, item todobus.TodoItem) error {
	const q = `
	UPDATE
		todo_items
	SET 
		description = :description,
		due_date = :due_date,
		file_id = :file_id,
		date_updated = :date_updated
	WHERE
		item_id = :item_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTodoItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a TodoItem from the database.
func (s *Store) Delete(ctx context.Context, item todobus.TodoItem) error {
	const q = `
	DELETE FROM
		todo_items
	WHERE
		item_id = :item_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTodoItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing TodoItems from the database.
func (s *Store) Query(ctx context.Context) ([]todobus.TodoItem, error) {
	const q = `
    SELECT
        item_id, description, due_date, file_id, date_created, date_updated
    FROM
        todo_items`

	buf := bytes.NewBufferString(q)

	var dbItems []dbTodoItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, buf, &dbItems); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusTodoItems(dbItems)
}

// QueryByID retrieves a specific TodoItem from the database by ID.
func (s *Store) QueryByID(ctx context.Context, itemID uuid.UUID) (todobus.TodoItem, error) {
	data := struct {
		ID string `db:"item_id"`
	}{
		ID: itemID.String(),
	}

	const q = `
	SELECT
		item_id, description, due_date, file_id, date_created, date_updated
	FROM
		todo_items
	WHERE 
		item_id = :item_id`

	var dbItem dbTodoItem
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbItem); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return todobus.TodoItem{}, fmt.Errorf("db: %w", todobus.ErrNotFound)
		}
		return todobus.TodoItem{}, fmt.Errorf("db: %w", err)
	}

	return toBusTodoItem(dbItem)
}
