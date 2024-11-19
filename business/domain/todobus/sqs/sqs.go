// Package sqs provides functionality for interacting with SQS message queues.
package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// Client manages interactions with SQS.
type Client struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string
}

// NewClient creates a new SQS client instance.
func NewClient(sqsClient sqsiface.SQSAPI, queueURL string) *Client {
	return &Client{
		sqsClient: sqsClient,
		queueURL:  queueURL,
	}
}

// SendMessage sends a message to the SQS queue.
func (c *Client) SendMessage(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = c.sqsClient.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}

// ReceiveMessages retrieves messages from the SQS queue.
func (c *Client) ReceiveMessages(ctx context.Context) ([]interface{}, error) {
	result, err := c.sqsClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: aws.Int64(10), // You can customize this number
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from SQS: %w", err)
	}

	messages := make([]interface{}, len(result.Messages))
	for i, msg := range result.Messages {
		var message interface{}
		if err := json.Unmarshal([]byte(*msg.Body), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message body: %w", err)
		}
		messages[i] = message
	}

	return messages, nil
}

// DeleteMessage removes a message from the SQS queue by its receipt handle.
func (c *Client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := c.sqsClient.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from SQS: %w", err)
	}

	return nil
}
