package inmemory_test

import (
	"context"
	"testing"

	"tekenajaqueue/queue/inmemory"
)

func TestQueue_Enqueue(t *testing.T) {
	q := inmemory.New()
	ctx := context.Background()
	t.Run("OK", func(t *testing.T) {
		err := q.Enqueue(ctx, "OK")
		if err != nil {
			t.Fatalf("scenario [%s] failed, got err:%s\n", t.Name(), err.Error())
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		err := q.Enqueue(ctx, "")
		if err == nil {
			t.Fatalf("scenario [%s] failed\n", t.Name())
		}
	})
}

func TestQueue_Dequeue(t *testing.T) {
	q := inmemory.New()
	ctx := context.Background()
	q.Enqueue(ctx, "OK")
	q.Enqueue(ctx, "OK_2")
	t.Run("message OK got pulled", func(t *testing.T) {
		msg, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("scenario [%s] failed, got err:%s\n", t.Name(), err.Error())
		}

		if msg != "OK" {
			t.Fatalf("scenario [%s] failed, pulling wrong message\n", t.Name())
		}
	})
}

func TestQueue_Clear(t *testing.T) {
	q := inmemory.New()
	ctx := context.Background()
	q.Enqueue(ctx, "OK")
	q.Enqueue(ctx, "OK_2")

	t.Run("OK", func(t *testing.T) {
		q.Clear()

		msg, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("scenario [%s] failed, got err:%s\n", t.Name(), err.Error())
		}

		if msg != "" {
			t.Fatalf("scenario [%s] failed, clearance queue failed\n", t.Name())
		}
	})
}
