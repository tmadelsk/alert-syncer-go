package rate

import "context"

type Limiter interface {
	Acquire(ctx context.Context) error
}

type NoopLimiter struct{}

func (l *NoopLimiter) Acquire(ctx context.Context) error {
	return nil
}
