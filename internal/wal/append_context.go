package wal

import "context"

func AppendWithContext(ctx context.Context, log *Log, record Record, syncWrite bool) error {
	return log.Append(record, syncWrite)
}
