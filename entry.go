package scheduler

import (
	"time"
)

type Entry struct {
	ID       string
	Index    int
	Next     time.Time
	Schedule Schedule
	Task     Task
}
