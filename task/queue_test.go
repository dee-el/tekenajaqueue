package task_test

import (
	"context"
	"testing"

	inmemory_queue "tekenajaqueue/queue/inmemory"

	task "tekenajaqueue/task"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	q := task.NewTaskQueue("test", inmemory_queue.New())
	err := q.Enqueue(ctx, &task.Task{ID: 1})
	if err != nil {
		t.Fatalf("err happened when Enqueue: %v", err.Error())
	}

	task, err := q.Dequeue(ctx)
	if err != nil {
		t.Fatalf("err happened when 1st pull: %v", err.Error())
	}

	if task == nil {
		t.Fatal("err happened when checking task in 1st pull: should not be empty")
	}

	task, err = q.Dequeue(ctx)
	if err != nil {
		t.Fatalf("err happened when 2nd pull: %v", err.Error())
	}

	if task != nil {
		t.Fatal("err happened when checking task in 2nd pull: should be empty")
	}
}
