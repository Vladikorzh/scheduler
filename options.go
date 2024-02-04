package scheduler

import "time"

type Options struct {
	Tick time.Duration
}

type Option interface {
	apply(opts *Options)
}

type OptionFunc func(opts *Options)

func (opt OptionFunc) apply(opts *Options) {
	opt(opts)
}

func Tick(tick time.Duration) Option {
	return OptionFunc(func(opts *Options) {
		opts.Tick = tick
	})
}
