package scheduler

import (
	"container/heap"
	"context"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Scheduler struct {
	tick    time.Duration
	heap    Heap
	running bool
	mx      sync.RWMutex
}

func New(options ...Option) *Scheduler {
	opts := Options{
		Tick: time.Second,
	}

	for _, opt := range options {
		opt.apply(&opts)
	}

	return &Scheduler{
		tick: opts.Tick,
	}
}

func (s *Scheduler) Schedule(schedule Schedule, task Task) {
	s.mx.Lock()
	defer s.mx.Unlock()

	heap.Push(&s.heap, Entry{
		ID:       uuid.New().String(),
		Index:    s.heap.Len(),
		Next:     schedule.Next(time.Now()),
		Schedule: schedule,
		Task:     task,
	})
}

func (s *Scheduler) ScheduleFunc(schedule Schedule, task TaskFunc) {
	s.Schedule(schedule, task)
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.mx.Lock()

	if s.running {
		s.mx.Unlock()

		return nil
	}

	s.running = true

	s.mx.Unlock()

	return s.run(ctx)
}

func (s *Scheduler) run(ctx context.Context) error {
	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-ticker.C:
			for {
				entry, ok := s.next(now)
				if !ok {
					break
				}

				s.dispatch(ctx, now, entry)
			}
		}
	}
}

func (s *Scheduler) next(now time.Time) (*Entry, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if len(s.heap) == 0 {
		return nil, false
	}

	if s.heap[0].Next.After(now) {
		return nil, false
	}

	return &s.heap[0], true
}

func (s *Scheduler) dispatch(ctx context.Context, now time.Time, entry *Entry) {
	go func() {
		entry.Task.Run(ctx)
	}()

	s.mx.Lock()
	defer s.mx.Unlock()

	if entry.Next = entry.Schedule.Next(now); entry.Next.After(now) {
		heap.Fix(&s.heap, entry.Index)
	} else {
		heap.Remove(&s.heap, entry.Index)
	}
}
