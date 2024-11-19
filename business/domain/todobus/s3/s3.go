// Package s3 provides functionality for interacting with S3 storage services.
package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// Client manages interactions with S3.
type Client struct {
	s3Client     s3iface.S3API
	s3Uploader   s3manageriface.UploaderAPI
	s3Downloader s3manageriface.DownloaderAPI
	bucketName   string
}

// NewClient creates a new S3 client instance.
func NewClient(s3Client s3iface.S3API, uploader s3manageriface.UploaderAPI, downloader s3manageriface.DownloaderAPI, bucketName string) *Client {
	return &Client{
		s3Client:     s3Client,
		s3Uploader:   uploader,
		s3Downloader: downloader,
		bucketName:   bucketName,
	}
}

// Upload uploads a file to S3 and returns its file ID (key).
func (c *Client) Upload(ctx context.Context, fileName string, data []byte) (string, error) {
	input := &s3manager.UploadInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileName),
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
	}

	result, err := c.s3Uploader.UploadWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return result.Location, nil
}

// Download retrieves a file from S3 by its file ID (key).
func (c *Client) Download(ctx context.Context, fileID string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer([]byte{})

	_, err := c.s3Downloader.DownloadWithContext(ctx, buff, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return buff.Bytes(), nil
}

// Delete removes a file from S3 by its file ID (key).
func (c *Client) Delete(ctx context.Context, fileID string) error {
	_, err := c.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
