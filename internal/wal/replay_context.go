package wal

import "context"

func ReplayWithContext(ctx context.Context, records []Record, apply func(Record) error) error {
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := apply(record); err != nil {
			return err
		}
	}
	return nil
}
