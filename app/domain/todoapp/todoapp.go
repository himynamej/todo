// Package todoapp provides functionality for todo application.
package todoapp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/himynamej/todo/app/sdk/errs"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/foundation/web"
)

type app struct {
	todoBus *todobus.Business
}

func newApp(todoBus *todobus.Business) *app {
	return &app{
		todoBus: todoBus,
	}
}

// CreateTodoItem handles creating a new TodoItem.
func (a *app) CreateTodoItem(ctx context.Context, r *http.Request) web.Encoder {
	var app NewTodoItem
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	// Use RFC3339 layout to parse the DueDate string
	parsedTime, err := time.Parse(time.RFC3339, app.DueDate)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("invalid date format: %w", err))
	}

	// Create the TodoItem using the business layer
	item, err := a.todoBus.Create(ctx, app.Description, parsedTime, nil, app.FileID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppTodoItem(item)
}

// UploadFile handles file uploads and stores them in an S3 bucket.
func (a *app) UploadFile(ctx context.Context, r *http.Request) web.Encoder {
	// Set a maximum size limit for file uploads (e.g., 10 MB)
	const maxFileSize = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(nil, r.Body, maxFileSize)

	// Extract the file name from a query parameter or custom header (adjust as necessary)
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		return errs.New(errs.InvalidArgument, fmt.Errorf("missing file name"))
	}

	// Read the file data directly from the request body
	fileData, err := io.ReadAll(r.Body)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("error reading file data: %w", err))
	}

	// Optional validation: Limit file size to a smaller maximum (e.g., 5 MB)
	if len(fileData) > 5<<20 { // 5 MB limit
		return errs.New(errs.InvalidArgument, fmt.Errorf("file size exceeds the maximum allowed limit of 5 MB"))
	}

	// Validate the file type by inspecting the initial bytes
	allowedFileTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"text/plain": true,
	}
	fileType := http.DetectContentType(fileData[:512])
	if !allowedFileTypes[fileType] {
		return errs.New(errs.InvalidArgument, fmt.Errorf("unsupported file type: %s", fileType))
	}

	// Upload the file to S3 using the business layer
	fileID, err := a.todoBus.UploadFile(ctx, fileName, fileData)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("error uploading file: %w", err))
	}

	// Return the file ID in the response using the FileUploadResponse model
	return FileUploadResponse{
		FileID: fileID,
	}
}

// DownloadFile handles downloading a file from S3.
func (a *app) DownloadFile(ctx context.Context, r *http.Request) web.Encoder {
	// Extract the file ID from URL parameters
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	// Retrieve the file data from S3 using the business layer
	fileData, err := a.todoBus.GetFile(ctx, qp.ID)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("error retrieving file: %w", err))
	}

	// Set a default file name for the downloaded file
	fileName := "downloaded-file"

	// Create a response to download the file
	return FileResponse{
		FileName:    fileName,
		ContentType: "application/octet-stream",
		Data:        fileData,
	}
}
