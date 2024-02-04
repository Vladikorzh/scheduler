package scheduler

import (
	"time"
)

type Schedule interface {
	Next(now time.Time) time.Time
}

type ScheduleFunc func(now time.Time) time.Time

func (fn ScheduleFunc) Next(now time.Time) time.Time {
	return fn(now)
}

type Every time.Duration

func (e Every) Next(now time.Time) time.Time {
	return now.Add(time.Duration(e) - time.Duration(now.Nanosecond())%time.Duration(e))
}

type At time.Time

func (at At) Next(_ time.Time) time.Time {
	return time.Time(at)
}
