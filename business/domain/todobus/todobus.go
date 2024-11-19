// Package todobus provides business access to TodoItem domain.
package todobus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/himynamej/todo/foundation/otel"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("todo item not found")
)

// Business manages the set of APIs for TodoItem access.
type Business struct {
	log      *logger.Logger
	storer   Storer
	sqsQueue SQSClient
	s3Client S3Client
}

// NewBusiness constructs a TodoItem business API for use.
func NewBusiness(log *logger.Logger, storer Storer, sqsQueue SQSClient, s3Client S3Client) *Business {
	return &Business{
		log:      log,
		storer:   storer,
		sqsQueue: sqsQueue,
		s3Client: s3Client,
	}
}

// Create adds a new TodoItem to the system, uploads the file to S3, and sends an SQS message.
func (b *Business) Create(ctx context.Context, description string, dueDate time.Time, fileData []byte, fileName string) (TodoItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.create")
	defer span.End()

	item := TodoItem{
		ID:          uuid.New(),
		Description: description,
		DueDate:     dueDate,
	}

	// Upload file to S3
	fileID, err := b.s3Client.Upload(ctx, fileName, fileData)
	if err != nil {
		return TodoItem{}, fmt.Errorf("s3 upload failed: %w", err)
	}
	item.FileID = fileID

	// Store item in the database
	if err := b.storer.Create(ctx, item); err != nil {
		return TodoItem{}, fmt.Errorf("create: %w", err)
	}

	// Send item data to SQS queue
	if err := b.sqsQueue.SendMessage(ctx, item); err != nil {
		b.log.Warn(ctx, "failed to send message to SQS", "error", err)
	}

	return item, nil
}

// Update modifies an existing TodoItem.
func (b *Business) Update(ctx context.Context, itemID uuid.UUID, newDescription string, newDueDate time.Time) (TodoItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.update")
	defer span.End()

	item, err := b.storer.QueryByID(ctx, itemID)
	if err != nil {
		return TodoItem{}, fmt.Errorf("update: %w", err)
	}

	item.Description = newDescription
	item.DueDate = newDueDate

	if err := b.storer.Update(ctx, item); err != nil {
		return TodoItem{}, fmt.Errorf("update: %w", err)
	}

	return item, nil
}

// Delete removes a specified TodoItem.
func (b *Business) Delete(ctx context.Context, itemID uuid.UUID) error {
	ctx, span := otel.AddSpan(ctx, "business.todobus.delete")
	defer span.End()

	item, err := b.storer.QueryByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if err := b.storer.Delete(ctx, item); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// Optionally delete the file from S3 if needed
	if err := b.s3Client.Delete(ctx, item.FileID); err != nil {
		b.log.Warn(ctx, "failed to delete file from S3", "fileID", item.FileID, "error", err)
	}

	return nil
}

// QueryByID finds a TodoItem by the specified ID.
func (b *Business) QueryByID(ctx context.Context, itemID uuid.UUID) (TodoItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.querybyid")
	defer span.End()

	item, err := b.storer.QueryByID(ctx, itemID)
	if err != nil {
		return TodoItem{}, fmt.Errorf("query: itemID[%s]: %w", itemID, err)
	}

	return item, nil
}

// Query retrieves a list of existing TodoItems.
func (b *Business) Query(ctx context.Context) ([]TodoItem, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.query")
	defer span.End()

	items, err := b.storer.Query(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return items, nil
}

// UploadFile uploads a file to S3 and returns the generated file ID.
func (b *Business) UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.uploadfile")
	defer span.End()

	// Validate inputs
	if fileName == "" || fileData == nil {
		return "", fmt.Errorf("uploadfile: invalid input: fileName[%s]", fileName)
	}

	fileID, err := b.s3Client.Upload(ctx, fileName, fileData)
	if err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	return fileID, nil
}

// GetFile retrieves a file from S3 by its file ID.
func (b *Business) GetFile(ctx context.Context, fileID string) ([]byte, error) {
	ctx, span := otel.AddSpan(ctx, "business.todobus.getfile")
	defer span.End()

	if fileID == "" {
		return nil, fmt.Errorf("getfile: fileID is required")
	}

	fileData, err := b.s3Client.Download(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("s3 download failed: %w", err)
	}

	return fileData, nil
}
