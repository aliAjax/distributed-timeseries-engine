package transport

import (
	"context"
	"time"
)

func BridgeDeadline(ctx context.Context) (time.Time, bool) {
	return ctx.Deadline()
}
