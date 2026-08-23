package wal

import "context"

func ScanWithContext(ctx context.Context, records []Record, visit func(int, Record) error) error {
	for index, record := range records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := visit(index, record); err != nil {
			return err
		}
	}
	return nil
}
