package task

import (
	"fmt"
	"time"
)

// Task contains information about Order fulfillment job
type Task struct {
	ID   int64     `json:"id"` // task id.
	Name string    `json:"name"`
	Due  time.Time `json:"due"` // minimum time for a Task to be executed.
}

func NewTask(ID int64, Name string, Due time.Time) *Task {
	return &Task{
		ID:   ID,
		Name: Name,
		Due:  Due,
	}
}

func (t *Task) Do() error {
	if !t.IsOverDue() {
		return fmt.Errorf("[%d] --> not yet time", t.ID)
	}

	fmt.Printf("[%d] --> %s executed\n", t.ID, t.Name)
	return nil
}

// IsOverDue checks if a task is already due, by calculating time.Now
func (t *Task) IsOverDue() bool {
	return time.Now().After(t.Due)
}
