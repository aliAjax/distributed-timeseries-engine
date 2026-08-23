package transport

import "context"

func RequestContextForLog(ctx context.Context) context.Context {
	return context.Background()
}
