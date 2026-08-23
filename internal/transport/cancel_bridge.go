package transport

import "context"

func BridgeCancellation(ctx context.Context, next func(context.Context) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return next(ctx)
}
