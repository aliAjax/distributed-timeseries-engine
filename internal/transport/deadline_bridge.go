package transport

import (
	"context"
	"time"
)

func BridgeDeadline(ctx context.Context) (time.Time, bool) {
	return context.Background().Deadline()
}
