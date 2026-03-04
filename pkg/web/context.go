package web

import (
	"context"
	"time"
)

// WithMinimumDeadline ensures ctx has at least min duration until deadline.
// Used by handlers and by main's importShield. Exported so main can use it.
func WithMinimumDeadline(ctx context.Context, min time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), min)
	}
	if deadline, ok := ctx.Deadline(); ok {
		if time.Until(deadline) >= min {
			return ctx, func() {}
		}
		return context.WithTimeout(context.WithoutCancel(ctx), min)
	}
	return context.WithTimeout(ctx, min)
}
