package wal

import "context"

func ReplayWithContext(ctx context.Context, records []Record, apply func(Record) error) error {
	for _, record := range records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := apply(record); err != nil {
			return err
		}
	}
	return nil
}
