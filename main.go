package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"tekenajaqueue/queue/inmemory"
	"tekenajaqueue/task"
	"tekenajaqueue/worker"
)

func handler(tq task.TaskQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := rand.Int63()
		tq.Enqueue(r.Context(), &task.Task{
			ID:   id,
			Name: fmt.Sprintf("hello_%d", id),
			Due:  time.Now().Add(5 * time.Second),
		})

		log.Printf("add queue: %d\n", id)

		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	ctx := context.Background()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	tq := task.NewTaskQueue("task_queue_1", inmemory.New())
	pool := worker.NewPool("pool_1", tq, worker.WithWorkerNum(100))

	// Set up HTTP server to add tasks
	http.HandleFunc("/tasks", handler(tq))

	go pool.Run(ctx)

	go func() {
		log.Printf("app server is up and running. Go to http://127.0.0.1:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting the HTTP server:", err)
		}
	}()

	<-quit
}
