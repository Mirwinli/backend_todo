package models

import "time"

type Task struct {
	Title       string
	Description string
	CreatedAt   time.Time
	DoneAt      *time.Time
	Duration    *time.Duration
	TaskId      int64
	IsDone      bool
}
