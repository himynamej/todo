package todobus

import (
	"context"

	"github.com/google/uuid"
)

// S3Client defines the interface for S3 operations.
type S3Client interface {
	Upload(ctx context.Context, fileName string, data []byte) (string, error)
	Download(ctx context.Context, fileID string) ([]byte, error)
	Delete(ctx context.Context, fileID string) error
}

// SQSClient defines the interface for SQS operations.
type SQSClient interface {
	SendMessage(ctx context.Context, message interface{}) error
	ReceiveMessages(ctx context.Context) ([]interface{}, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, item TodoItem) error
	Update(ctx context.Context, item TodoItem) error
	Delete(ctx context.Context, item TodoItem) error
	QueryByID(ctx context.Context, itemID uuid.UUID) (TodoItem, error)
	Query(ctx context.Context) ([]TodoItem, error)
}
