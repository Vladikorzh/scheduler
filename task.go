package scheduler

import "context"

type Task interface {
	Run(ctx context.Context)
}

type TaskFunc func(ctx context.Context)

func (fn TaskFunc) Run(ctx context.Context) {
	fn(ctx)
}
