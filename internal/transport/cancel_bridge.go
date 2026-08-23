package transport

import "context"

func BridgeCancellation(ctx context.Context, next func(context.Context) error) error {
	return next(context.Background())
}
