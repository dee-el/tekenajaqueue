package inmemory

import (
	"container/list"
	"context"
	"errors"
	"sync"
)

// Queue implements queue.Queue using linked list as storage
type Queue struct {
	list *list.List
	mut  sync.Mutex
}

// New returns new Queue
func New() *Queue {
	return &Queue{
		list: list.New(),
	}
}

// threadSafe ensure fn runs in synchronous manner
func (q *Queue) threadSafe(fn func()) {
	q.mut.Lock()
	defer q.mut.Unlock()
	fn()
}

// Enqueue stores data to queue by
// adding the data to the tail of underlying linked list
func (q *Queue) Enqueue(ctx context.Context, msg string) error {
	if msg == "" {
		return errors.New("empty message")
	}

	q.threadSafe(func() {
		q.list.PushBack(msg)
	})
	return nil
}

// Dequeue removes next data from queue by returning
// data at the head of underlying linked list
// and returns the removed data
func (q *Queue) Dequeue(ctx context.Context) (string, error) {
	var msg string
	q.threadSafe(func() {
		if q.list.Len() == 0 {
			return
		}

		elem := q.list.Front()
		if d, ok := elem.Value.(string); ok && d != "" {
			q.list.Remove(elem)
			msg = d
		}
	})

	return msg, nil
}

// Clear removes all data in the queue by
// removing all underlying linked list elements
func (q *Queue) Clear() error {
	curr := q.list.Front()
	for {
		if curr == nil {
			break
		}

		next := curr.Next()
		q.list.Remove(curr)
		curr = next
	}

	return nil
}
