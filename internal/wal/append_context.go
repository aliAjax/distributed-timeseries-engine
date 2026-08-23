package wal

import "context"

func AppendWithContext(ctx context.Context, log *Log, record Record, syncWrite bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return log.Append(record, syncWrite)
}
