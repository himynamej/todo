package todobus_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/business/domain/todobus/mocks"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/himynamej/todo/foundation/otel"
)

func BenchmarkInsertTodoItem(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}
	mockLogger := logger.New(os.Stdout, logger.LevelInfo, "SALES", traceIDFn)

	mockStorer := mocks.NewMockStorer(ctrl)
	mockSQSClient := mocks.NewMockSQSClient(ctrl)
	mockS3Client := mocks.NewMockS3Client(ctrl)

	bus := todobus.NewBusiness(mockLogger, mockStorer, mockSQSClient, mockS3Client)

	// Create a sample TodoItem.
	description := "Sample TodoItem"
	dueDate := time.Now().Add(24 * time.Hour)
	fileData := []byte("Sample file data")
	fileName := "sample.txt"

	// Mock the expected interactions
	mockStorer.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockS3Client.EXPECT().Upload(gomock.Any(), fileName, fileData).Return("mock-file-id", nil).AnyTimes()
	mockSQSClient.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Create(context.Background(), description, dueDate, fileData, fileName)
		if err != nil {
			b.Fatalf("failed to insert TodoItem: %v", err)
		}
	}
}

func BenchmarkUploadFileToS3(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockS3Client := mocks.NewMockS3Client(ctrl)

	fileData := []byte("Sample file data for benchmarking")
	fileName := "benchmark-file.txt"

	// Mock the expected interactions
	mockS3Client.EXPECT().Upload(gomock.Any(), fileName, fileData).Return("mock-file-id", nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockS3Client.Upload(context.Background(), fileName, fileData)
		if err != nil {
			b.Fatalf("failed to upload file to S3: %v", err)
		}
	}
}

func BenchmarkSendMessageToSQS(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockSQSClient := mocks.NewMockSQSClient(ctrl)

	message := map[string]string{
		"ID":          uuid.New().String(),
		"Description": "Sample TodoItem for SQS",
	}

	// Mock the expected interactions
	mockSQSClient.EXPECT().SendMessage(gomock.Any(), message).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := mockSQSClient.SendMessage(context.Background(), message)
		if err != nil {
			b.Fatalf("failed to send message to SQS: %v", err)
		}
	}
}
