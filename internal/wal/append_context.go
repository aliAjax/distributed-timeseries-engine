package wal

import "context"

func AppendWithContext(ctx context.Context, log *Log, record Record, syncWrite bool) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return log.Append(record, syncWrite)
}
