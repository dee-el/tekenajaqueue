// Package queue provides an interface for a simple message queue.
package queue

import "context"

// Queue is an interface representing a message queue.
type Queue interface {
	// Enqueue adds a message to the queue.
	// It takes a context and the message to be enqueued.
	// Returns an error if the operation fails.
	Enqueue(ctx context.Context, msg string) error

	// Dequeue retrieves a message from the queue.
	// It takes a context and returns the dequeued message.
	// Returns an error if the operation fails.
	Dequeue(ctx context.Context) (string, error)

	// Clear removes all messages from the queue.
	// Returns an error if the operation fails.
	Clear() error
}
