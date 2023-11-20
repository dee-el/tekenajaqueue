package task

import (
	"context"
	"encoding/json"

	"tekenajaqueue/queue"
)

// TaskQueue define task queue operations
type TaskQueue interface {
	Enqueue(ctx context.Context, t *Task) error
	Dequeue(ctx context.Context) (*Task, error)
	Clear() error
}

// Queue an standard implementation for task.Queue.
// The responsibilities is transform Task to message on queue or vice versa.
type taskQueue struct {
	code    string
	queue   queue.Queue
	encoder Encoder
	decoder Decoder
}

// New by default task is encoded - decoded via JSON
func NewTaskQueue(code string, q queue.Queue) *taskQueue {
	return NewTaskQueueWithOption(code, q, JSONEncodeDecodeOption())
}

func NewTaskQueueWithOption(code string, q queue.Queue, opts ...QueueOption) *taskQueue {
	queue := &taskQueue{
		code:  code,
		queue: q,
	}

	for _, opt := range opts {
		opt(queue)
	}

	return queue
}

// QueueOption is function to configure created queue
type QueueOption func(q *taskQueue)

func JSONEncodeDecodeOption() QueueOption {
	return func(q *taskQueue) {
		q.encoder = JSONEncoder
		q.decoder = JSONDecoder
	}
}

// Encoder transform Task to message string
type Encoder func(t *Task) (string, error)

// JSONEncoder is JSON implementation of encoder
var JSONEncoder = func(t *Task) (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(b), err
}

// Decoder transform message string to Task
type Decoder func(message string) (*Task, error)

var JSONDecoder = func(message string) (*Task, error) {
	var t *Task
	err := json.Unmarshal([]byte(message), &t)
	if err != nil {
		return nil, err
	}

	return t, err
}

func (q *taskQueue) Enqueue(ctx context.Context, t *Task) error {
	msg, err := q.encoder(t)
	if err != nil {
		return err
	}

	return q.queue.Enqueue(ctx, msg)
}

// Dequeue dequeue data from queue and returns it as *Task
func (q *taskQueue) Dequeue(ctx context.Context) (*Task, error) {
	msg, err := q.queue.Dequeue(ctx)
	if err != nil {
		return nil, err
	}

	if msg == "" {
		return nil, nil
	}

	return q.decoder(msg)
}

func (q *taskQueue) Clear() error {
	return q.queue.Clear()
}
