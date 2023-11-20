// Package worker provides a simple worker pool for processing tasks from a task queue.
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"tekenajaqueue/task"
)

// Pool is the interface for a worker pool.
type Pool interface {
	Run(ctx context.Context)
}

// pool is an implementation of the Pool interface.
type pool struct {
	name    string
	queue   task.TaskQueue
	workers chan struct{}
}

// option is a functional option for configuring the worker pool.
type option func(*pool)

// WithWorkerNum sets the number of workers in the pool.
func WithWorkerNum(num int) option {
	return func(p *pool) {
		p.workers = make(chan struct{}, num)
	}
}

// NewPool creates a new worker pool with the specified name and task queue.
// It accepts optional configurations via functional options.
func NewPool(name string, queue task.TaskQueue, opts ...option) *pool {
	p := &pool{
		name:    name,
		queue:   queue,
		workers: make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Run starts the worker pool, continuously dequeuing tasks and processing them.
// It uses a ticker to control the rate at which workers are allowed to pick up new tasks.
// The pool can be stopped by canceling the context.
func (d *pool) Run(ctx context.Context) {
	t := time.NewTicker(time.Duration(100) * time.Millisecond)
	defer t.Stop()

	for {
		task, err := d.queue.Dequeue(ctx)
		if err != nil {
			log.Println(fmt.Printf("[%s] failed to dequeue: %v", d.name, err))
		}

		select {
		case <-ctx.Done():
			// end progress
			return
		case <-t.C:
			d.workers <- struct{}{}
		}

		go func() {
			// release the semaphore
			defer func() {
				<-d.workers
			}()

			// Do work based on the message, for example:
			if task == nil {
				return
			}

			err := task.Do()
			if err != nil {
				log.Printf("[%s] failed to process: %s \n", d.name, err.Error())

				// return it to queue, for a real production project should be redirected into DLQ
				d.queue.Enqueue(ctx, task)
				return
			}

		}()
	}
}
